package repositories

import (
	"auth-api/internal/models"
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var adminCollection *mongo.Collection
var adminSessionCollection *mongo.Collection

// initialises the admin user collection
func SetAdminCollection(collection *mongo.Collection) {
	adminCollection = collection
}

func SetAdminSessionCollection(collection *mongo.Collection) {
	adminSessionCollection = collection
}

func SaveAdmin(adminUser models.AdminUser) (models.AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := adminCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return models.AdminUser{}, err
	}

	adminUser.ID = result.InsertedID.(primitive.ObjectID)
	return adminUser, nil
}

func GetAllAdmins() ([]models.AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := adminCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var admins []models.AdminUser
	if err := cursor.All(ctx, &admins); err != nil {
		return nil, err
	}

	return admins, nil
}

func CreateAdminSessionIndexes() error {
	_, err := adminSessionCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "access_token", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "access_expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(int32(24 * 60 * 60)),
		},
	})
	return errors.Wrap(err, "failed to create session indexes")
}

func SaveAdminSession(adminSession models.AdminSession) (models.AdminSession, error) {
	adminSession.CreatedAt = time.Now()
	result, err := adminSessionCollection.InsertOne(context.Background(), adminSession)
	if err != nil {
		return models.AdminSession{}, err
	}
	adminSession.AdminID = result.InsertedID.(primitive.ObjectID)
	return adminSession, nil
}

func GetAdminSessionByAccessToken(adminAccessToken string) (models.AdminSession, error) {
	var adminSession models.AdminSession
	err := adminSessionCollection.FindOne(context.Background(), bson.M{"access_token": adminAccessToken}).Decode(&adminSession)
	return adminSession, errors.Wrap(err, "failed to get session by access token")
}

func GetActiveSessionsByAdminID(adminID primitive.ObjectID) ([]models.AdminSession, error) {
	filter := bson.M{
		"admin_id":         adminID,
		"accss_expires_at": bson.M{"$gt": time.Now()},
	}

	cursor, err := adminSessionCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active sessions")
	}

	var adminSessions []models.AdminSession
	if err = cursor.All(context.Background(), &adminSessions); err != nil {
		return nil, errors.Wrap(err, "failed to decode sessions")
	}

	return adminSessions, nil
}

func DeleteAdminSession(id primitive.ObjectID) error {
	_, err := adminSessionCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func UpdateAdminSession(adminSession models.AdminSession) (models.AdminSession, error) {
	filter := bson.M{"_id": adminSession.AdminID}
	update := bson.M{"$set": adminSession}

	_, err := adminSessionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return models.AdminSession{}, errors.Wrap(err, "failed to update admin session")
	}
	return adminSession, nil
}

func InvalidateAdminSessionByAccessToken(accessToken string) error {
	filter := bson.M{"access_token": accessToken}
	update := bson.M{"$set": bson.M{"access_expires_at": time.Now()}}

	_, err := adminSessionCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to invalidate admin session by access token")
	}
	return nil
}

func GetAdminByUsername(username string) (models.AdminUser, error) {
	var adminUser models.AdminUser
	err := adminCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&adminUser)
	return adminUser, errors.Wrap(err, "failed to get admin by username")
}
