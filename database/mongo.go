package database

import (
	"context"
	"errors"

	"github.com/MRegterschot/GoController/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var dbName string

func GetCollection(name string) *mongo.Collection {
	return mongoClient.Database(dbName).Collection(name)
}

func Connect() error {
	uri := config.AppEnv.MongoUri
	if uri == "" {
		return errors.New("you must set your 'MONGO_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	database := config.AppEnv.MongoDb
	if database == "" {
		return errors.New("you must set your 'MONGO_DB' environmental variable")
	} else {
		dbName = database
	}

	var err error
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return errors.New("can't verify a connection")
	}

	return nil
}

func Disconnect() {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}
