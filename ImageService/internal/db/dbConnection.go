package db

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var MongoClient *mongo.Client

func ConnectDB() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Не удалось загрузить файл .env: %v", err)
	}

	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	connectionStr := fmt.Sprintf("mongodb://%s:%s@localhost:27017/?authSource=admin", username, password)

	clientOptions := options.Client().ApplyURI(connectionStr)

	MongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Не удалось подключиться к MongoDB: %v", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer pingCancel()

	err = MongoClient.Ping(pingCtx, nil)
	if err != nil {
		log.Fatalf("Не удалось выполнить ping к MongoDB: %v", err)
	}

	log.Println("Успешное подключение к MongoDB")
}
