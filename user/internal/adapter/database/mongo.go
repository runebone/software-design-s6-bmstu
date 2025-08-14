package database

import (
	"context"
	"fmt"
	"log"
	"user/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg config.MongoConfig) (*mongo.Database, error) {
	URI := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(cfg.DBName)

	log.Println("Connected to MongoDB successfully")
	return db, nil
}
