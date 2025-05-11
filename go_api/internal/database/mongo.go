package database

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	databaseInstance *mongo.Database
	mongoOnce        sync.Once
)

func ConnectMongo(uri string, databaseName string) *mongo.Database {
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		clientInstance, err := mongo.Connect(options.Client().ApplyURI(uri))

		if err != nil {
			panic(err)
		}

		if err := clientInstance.Ping(ctx, nil); err != nil {
			panic(err)
		}
		databaseInstance = clientInstance.Database(databaseName)
	})

	return databaseInstance
}

func GetCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
