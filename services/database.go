package services

import (
	"context"
	"log"
	"time"

	"link-shortener/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const(
	mongoURI = "mongodb+srv://atharv24:JYMGr9h7jkdCF8fw@cluster0.mskfpan.mongodb.net/?retryWrites=true&w=majority"
	database = "link-shortener"
)

var(
	db Database
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectDB() (error) {
	// Connect to MongoDB database
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	db = Database{
		Client:   client,
		Database: client.Database(database),
	}
	log.Println("Connected to DB")
	return nil
}

func DisconnectDB() {
	// Close MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.Client.Disconnect(ctx)
}

func AddLink(l *models.Link) error {
	// Insert new link into database
	_, err := db.Database.Collection("links").InsertOne(context.TODO(), l)
	if err != nil {
		return err
	}
	return nil
}

func GetLinkByShortURL(shortURL string) (*models.Link, error) {
	// Retrieve link from database
	var l models.Link
	err := db.Database.Collection("links").FindOne(context.TODO(), bson.M{"short_url": shortURL}).Decode(&l)
	log.Println(l, err)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &l, nil
}

func DeleteExpiredLinks(db *Database) error {
	// Delete expired links from database
	_, err := db.Database.Collection("links").DeleteMany(context.TODO(), bson.M{
		"created_at": bson.M{
			"$lt": time.Now().Add(-30 * 24 * time.Hour),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
