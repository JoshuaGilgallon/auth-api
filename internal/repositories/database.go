package repositories

import (
    "context"
    "auth-api/internal/config"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitDatabase(config *internalconfig.DatabaseConfig) error {
    ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
    defer cancel()

    clientOptions := options.Client().
        ApplyURI(config.URI).
        SetServerSelectionTimeout(config.ConnectTimeout).
        SetConnectTimeout(config.ConnectTimeout)

    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        return err
    }

    // Ping the database
    if err := client.Ping(ctx, nil); err != nil {
        return err
    }

    // Initialize collections
    db := client.Database(config.DatabaseName)
    if err := SetSessionCollection(db.Collection("sessions")); err != nil {
        return err
    }
    
    SetUserCollection(db.Collection("users"))
    return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
    if client != nil {
        return client.Disconnect(context.Background())
    }
    return nil
}
