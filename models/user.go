package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

type User struct {
	Email     string `json:"email" bson:"email,omitempty" validate:"required,email"`
	LastName  string `json:"last_name" bson:"last_name,omitempty" validate:"required,gte=2,lte=40"`
	Country   string `json:"country" bson:"country,omitempty" validate:"required,gte=3,lte=20"`
	City      string `json:"city" bson:"city,omitempty" validate:"required,gte=3,lte=30"`
	Gender    string `json:"gender" bson:"gender,omitempty" validate:"required,gte=0,lte=10"`
	BirthDate string `json:"birth_date" bson:"birth_date,omitempty" validate:"required,gte=0,lte=100"`
}

var validate *validator.Validate

// Users() returns all users which found in collection
func Users(page int) []User {
	var users []User
	collection := mongoConnect()
	opts := options.Find().SetSkip(int64(30 * (page - 1))).SetLimit(30)
	cursor, _ := collection.Find(context.TODO(), bson.D{}, opts)
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		checkErr(err)
		users = append(users, user)
	}
	return users
}
func (u *User) Create(user User) (*User, error) {
	collection := mongoConnect()
	validate = validator.New()

	// creating unique index for email field
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"email", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		return nil, errors.New("email already exists")
	}
	err = validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("field %s is required", err.Field())
		}
	}
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (u *User) Update(id primitive.ObjectID, user User) (User, error) {
	client := mongoConnect()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.D{
		{"email", user.Email},
		{"last_name", user.LastName},
		{"country", user.Country},
		{"city", user.City},
		{"gender", user.Gender},
		{"birth_date", user.BirthDate},
	}}
	_, err := client.UpdateOne(context.TODO(), filter, update)
	checkErr(err)
	return user, nil
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// this is helper function to connect to database
func mongoConnect() *mongo.Collection {
	client, err := mongo.Connect(context.TODO(), options.Client())
	checkErr(err)
	collection := client.Database("test").Collection("users")

	return collection
}

func test() {
	coll := mongoConnect()
	_, err := coll.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"email", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		})
	checkErr(err)
}
