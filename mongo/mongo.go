package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB_Client struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type userResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Level        int64              `bson:"level,omitempty"`
	Token        string             `bson:"token,omitempty"`
	Password     string             `bson:"password,omitempty"`
	Email        string             `bson:"email,omitempty"`
	RegisterDate float64            `bson:"reg_time,omitempty"`
	Payment      string             `bson:"payment_method,omitempty"`
}

func NewMongoDB(url string) (*MongoDB_Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	collection := client.Database("users").Collection("users")

	return &MongoDB_Client{
		Client:     client,
		Collection: collection,
	}, err
}

func (m *MongoDB_Client) RefreshToken(oldToken string, newToken string) {
	filtercursor, err := m.Collection.Find(context.TODO(), bson.M{"token": oldToken})
	if err != nil {
		panic(err)
	}
	var result []userResult
	if err = filtercursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
		return
	}
	if len(result) == 0 {
		return
	}
	updatedUser := userResult{
		Token:        newToken,
		ID:           result[0].ID,
		Name:         result[0].Name,
		Level:        result[0].Level,
		Password:     result[0].Password,
		Email:        result[0].Email,
		RegisterDate: result[0].RegisterDate,
		Payment:      result[0].Payment,
	}
	m.Collection.UpdateOne(context.TODO(), bson.M{"token": oldToken}, bson.M{"$set": updatedUser})
	if err != nil {
		panic(err)
	}

}
