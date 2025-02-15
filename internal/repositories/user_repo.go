package repositories

import (
	"auth-api/internal/models"
	"context"
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
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}
