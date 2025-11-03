package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client 

func ConnectDB() *mongo.Client {
	dsn := os.Getenv("MONGO_DSN")
	if dsn == "" {
		dsn = "mongodb://localhost:27017/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalf("Gagal membuat klien MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Koneksi database MongoDB gagal: %v", err)
	}

	log.Println("Koneksi ke MongoDB berhasil!")
	DB = client

	return client
}

func GetDB() *mongo.Client {
	return DB
}

func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
    
	return DB.Ping(ctx, nil) 
}

func CloseDB(client *mongo.Client) {
    if client == nil {
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
    
    if err := client.Disconnect(ctx); err != nil {
        log.Fatalf("Gagal menutup koneksi MongoDB: %v", err)
    }
    log.Println("Koneksi MongoDB ditutup.")
}