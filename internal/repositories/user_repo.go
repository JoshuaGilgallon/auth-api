package repositories

import (
	"auth-api/internal/models"
	internalconfig "auth-api/internal/config"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client *mongo.Client
var userCollection *mongo.Collection

func InitDatabase(cfg *internalconfig.DatabaseConfig) error {
	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetServerSelectionTimeout(cfg.ConnectTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()
	
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if (err != nil) {
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if (err != nil) {
		return err
	}

	userCollection = client.Database(cfg.DatabaseName).Collection("users")
	return nil
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
