package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var statsCollection *mongo.Collection

// initialises the admin user collection
func SetStatsCollection(collection *mongo.Collection) {
	statsCollection = collection
}

// checks and initializes the login stats document
func InitLoginStats() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"type": "logins"}
	update := bson.M{
		"$setOnInsert": bson.M{
			"type":  "logins",
			"value": bson.M{},
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := statsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	_, err = statsCollection.UpdateOne(ctx, filter, update)
	return err
}

func IncreaseLoginCount() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentWeek := time.Now().Format("2006-01-02")
	filter := bson.M{"type": "logins"}
	update := bson.M{
		"$inc": bson.M{
			"value." + currentWeek: 1,
		},
	}

	_, err := statsCollection.UpdateOne(ctx, filter, update)
	return err
}

func GetLoginCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentWeek := time.Now().Format("2006-01-02")
	filter := bson.M{"type": "logins"}
	projection := bson.M{
		"value." + currentWeek: 1,
		"_id":                  0,
	}

	var result struct {
		Value map[string]int `bson:"value"`
	}

	err := statsCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return result.Value[currentWeek], nil
}

func GetLastWeekLoginCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lastWeek := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	filter := bson.M{"type": "logins"}
	projection := bson.M{
		"value." + lastWeek: 1,
		"_id":               0,
	}

	var result struct {
		Value map[string]int `bson:"value"`
	}

	err := statsCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return result.Value[lastWeek], nil
}

func GetWeeklyStats() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentWeekLogins, err := GetLoginCount()
	if err != nil {
		return nil, err
	}

	lastWeekLogins, err := GetLastWeekLoginCount()
	if err != nil {
		return nil, err
	}

	var loginPercentageIncrease int
	if lastWeekLogins > 0 {
		loginPercentageIncrease = int(float64(currentWeekLogins-lastWeekLogins) / float64(lastWeekLogins) * 100)
	} else {
		loginPercentageIncrease = 100
	}

	usersCollection := statsCollection.Database().Collection("users")
	filter := bson.M{
		"created_at": bson.M{
			"$gte": time.Now().AddDate(0, 0, -7),
			"$lt":  time.Now(),
		},
	}
	currentWeekUsers, err := usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	filter = bson.M{
		"created_at": bson.M{
			"$gte": time.Now().AddDate(0, 0, -14),
			"$lt":  time.Now().AddDate(0, 0, -7),
		},
	}
	lastWeekUsers, err := usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	var userPercentageIncrease int
	if lastWeekUsers > 0 {
		userPercentageIncrease = int(float64(currentWeekUsers-lastWeekUsers) / float64(lastWeekUsers) * 100)
	} else {
		userPercentageIncrease = 100
	}

	return map[string]interface{}{
		"logins_change": loginPercentageIncrease,
		"users_change":  userPercentageIncrease,
	}, nil
}
