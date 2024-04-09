package database

import (
	"context"
	"fmt"
	"log"
	"password-manager/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient   *mongo.Client
	collection *mongo.Collection
)

func Connect(uri string) error {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Connected to MongoDB!")

	dbClient = client
	collection = client.Database("Password-Manager").Collection("Users")

	return nil
}

func GetClient() *mongo.Client {
	return dbClient
}

func CreateUser(user models.User) error {
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Failed to create user:", err)
	}
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	return user, err
}
