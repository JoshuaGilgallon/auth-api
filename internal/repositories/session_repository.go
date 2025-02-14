package repositories

import (
	"auth-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/pkg/errors"
)

var sessionCollection *mongo.Collection

// SetSessionCollection initializes the session collection and creates necessary indexes
func SetSessionCollection(collection *mongo.Collection) error {
	sessionCollection = collection
	return createSessionIndexes()
}

func createSessionIndexes() error {
	_, err := sessionCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "access_token", Value: 1},
				{Key: "refresh_token", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "refresh_expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(int32(24 * 60 * 60)),
		},
	})
	return errors.Wrap(err, "failed to create session indexes")
}

func SaveSession(session models.Session) (models.Session, error) {
	session.CreatedAt = time.Now()
	result, err := sessionCollection.InsertOne(context.Background(), session)
	if err != nil {
		return models.Session{}, err
	}
	session.ID = result.InsertedID.(primitive.ObjectID)
	return session, nil
}

func GetSessionByAccessToken(accessToken string) (models.Session, error) {
	var session models.Session
	err := sessionCollection.FindOne(context.Background(), bson.M{"access_token": accessToken}).Decode(&session)
	return session, errors.Wrap(err, "failed to get session by access token")
}

func GetSessionByRefreshToken(refreshToken string) (models.Session, error) {
	var session models.Session
	err := sessionCollection.FindOne(context.Background(), bson.M{"refresh_token": refreshToken}).Decode(&session)
	return session, errors.Wrap(err, "failed to get session by refresh token")
}

func DeleteSession(id primitive.ObjectID) error {
	_, err := sessionCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func GetActiveSessionsByUserID(userID primitive.ObjectID) ([]models.Session, error) {
	filter := bson.M{
		"user_id": userID,
		"refresh_expires_at": bson.M{"$gt": time.Now()},
	}
	
	cursor, err := sessionCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active sessions")
	}
	
	var sessions []models.Session
	if err = cursor.All(context.Background(), &sessions); err != nil {
		return nil, errors.Wrap(err, "failed to decode sessions")
	}
	
	return sessions, nil
}

func UpdateSession(session models.Session) (models.Session, error) {
	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": session}
	
	_, err := sessionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return models.Session{}, errors.Wrap(err, "failed to update session")
	}
	return session, nil
}
