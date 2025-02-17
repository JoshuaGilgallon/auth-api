package repositories

import (
	"auth-api/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var adminCollection *mongo.Collection

// initialises the user collection
func SetAdminCollection(collection *mongo.Collection) {
	adminCollection = collection
}

func SaveAdmin(adminUser models.AdminUser) (models.AdminUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return models.AdminUser{}, err
	}

	adminUser.ID = result.InsertedID.(primitive.ObjectID)
	return adminUser, nil
}
