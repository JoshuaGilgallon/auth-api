package repositories

import (
	"auth-api/internal/models"
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

// initialises the user collection
func SetUserCollection(collection *mongo.Collection) {
	userCollection = collection
}

func SaveUser(user models.User) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the user has an ID, update the existing document
	if !user.ID.IsZero() {
		filter := bson.M{"_id": user.ID}
		update := bson.M{"$set": user}

		_, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return models.User{}, fmt.Errorf("failed to update user: %w", err)
		}

		return user, nil
	}

	// For new users, insert a new document
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user: %w", err)
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
			log.Printf("user not found with id: %s", id)
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"email": email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("user not found with email: %s", email)
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil
}

func GetUsersByTimeCreatedRange(startTime, endTime time.Time, skip int64, limit int64) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	// Get total count first
	total, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results
	options := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := userCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUsersByTimeUpdatedRange(startTime, endTime time.Time, skip int64, limit int64) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"updated_at": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	// Get total count first
	total, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results
	options := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := userCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetTotalUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func SearchUsersByFields(criteria models.UserAdvancedSearchCriteria) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orConditions := []bson.M{}

	// Add name search conditions if provided
	if criteria.FirstName != "" {
		orConditions = append(orConditions, bson.M{
			"first_name": bson.M{
				"$regex":   "(?i)" + regexp.QuoteMeta(criteria.FirstName),
				"$options": "i",
			},
		})
	}

	if criteria.LastName != "" {
		orConditions = append(orConditions, bson.M{
			"last_name": bson.M{
				"$regex":   "(?i)" + regexp.QuoteMeta(criteria.LastName),
				"$options": "i",
			},
		})
	}

	if criteria.Email != "" {
		emailConditions := []bson.M{
			{
				"email": bson.M{
					"$regex":   "(?i)" + regexp.QuoteMeta(criteria.Email),
					"$options": "i",
				},
			},
		}
		orConditions = append(orConditions, bson.M{"$or": emailConditions})
	}

	filter := bson.M{}
	if len(orConditions) > 0 {
		filter = bson.M{"$or": orConditions}
	}

	// First get total count
	total, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip value
	skip := (criteria.PageNumber - 1) * criteria.PageSize

	// Add pagination options
	options := options.Find().
		SetSkip(skip).
		SetLimit(criteria.PageSize)

	// Execute the query
	cursor, err := userCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Use a map to deduplicate users by ID
	userMap := make(map[string]models.User)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, total, err
		}
		// Use ID as key for deduplication
		userMap[user.ID.Hex()] = user
	}

	err = cursor.Err()
	if err != nil {
		return nil, total, err
	}

	// Convert map back to slice
	users := make([]models.User, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, user)
	}

	return users, total, nil
}

func SimpleSearchUsers(searchTerm string, skip int64, limit int64) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create search conditions for all fields
	orConditions := []bson.M{
		// Name searches (case insensitive)
		{
			"first_name": bson.M{
				"$regex":   "(?i)" + regexp.QuoteMeta(searchTerm),
				"$options": "i",
			},
		},
		{
			"last_name": bson.M{
				"$regex":   "(?i)" + regexp.QuoteMeta(searchTerm),
				"$options": "i",
			},
		},
		// Email conditions
		{
			"$or": []bson.M{
				{
					"email": bson.M{
						"$regex":   "(?i)" + regexp.QuoteMeta(searchTerm),
						"$options": "i",
					},
				},
			},
		},
	}

	filter := bson.M{"$or": orConditions}

	// First, get total count without pagination
	total, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Then get paginated results
	options := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := userCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Use map to deduplicate results
	userMap := make(map[string]models.User)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		userMap[user.ID.Hex()] = user
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	// Convert map to slice
	users := make([]models.User, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, user)
	}

	return users, total, nil
}
