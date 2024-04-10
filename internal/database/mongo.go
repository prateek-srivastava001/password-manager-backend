package database

import (
	"context"
	"errors"
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

func AddCredential(email string, credential models.Credential) error {
	filter := bson.M{"email": email}

	update := bson.M{
		"$push": bson.M{"credentials": credential},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("some error occured")
	}

	return nil
}

func GetCredentialsForUser(email string) ([]models.Credential, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Credentials, nil
}

func EditCredential(email string, credentialID string, updatedCredential models.Credential) error {
	filter := bson.M{"email": email, "credentials.id": credentialID}

	update := bson.M{
		"$set": bson.M{
			"credentials.$.email":    updatedCredential.Email,
			"credentials.$.password": updatedCredential.Password,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("credential not found or no update required")
	}

	return nil
}

func DeleteCredential(email string, credentialID string) error {
	filter := bson.M{"email": email}
	update := bson.M{
		"$pull": bson.M{"credentials": bson.M{"id": credentialID}},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("credential not found or no deletion required")
	}

	return nil
}
