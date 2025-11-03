package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() *mongo.Client {
	connURI := "mongodb://localhost:27017/"
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Gagal membuat klien MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping database MongoDB: %v", err)
	}

	log.Println("Koneksi ke MongoDB berhasil!")
	return client
}

func LoggerMiddleware(c *fiber.Ctx) error {
	fmt.Println("Request:", c.Method(), c.Path())
	return c.Next()
}