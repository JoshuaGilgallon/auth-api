package repositories

import (
	"auth-api/internal/models"
	"auth-api/internal/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

// initialises the user collection
func SetUserCollection(collection *mongo.Collection) {
	userCollection = collection
}

func SaveUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func GetUserByID(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	filter := bson.M{"_id": objID}
	err = userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found with id: %s", id)
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashed_email := utils.HashSHA(email)

	var user models.User
	filter := bson.M{"email_hash": hashed_email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found with email: %s", email)
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUserByPhoneNumber(phoneNumber string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashed_phone_number := utils.HashSHA(phoneNumber)

	var user models.User
	filter := bson.M{"phone_number_hash": hashed_phone_number}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found with phone number: %s", phoneNumber)
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUsersByTimeCreatedRange(startTime, endTime time.Time) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersByTimeUpdatedRange(startTime, endTime time.Time) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"updated_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	cursor, err := userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
