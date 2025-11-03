package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Uploads struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UploadsName  string             `json:"Uploads_name" bson:"Uploads_name"`
	OriginalName string             `json:"original_name"bson:"original_name"`
	UploadsPath  string             `json:"Uploads_path" bson:"Uploads_path"`
	UploadsSize  int64              `json:"Uploads_size" bson:"Uploads_size"`
	UploadsType  string             `json:"Uploads_type" bson:"Uploads_type"`
	UploadedAt   time.Time          `json:"uploaded_at" bson:"uploaded_at"`
}

type UploadsResponse struct {
	ID           string    `json:"id"`
	UploadsName  string    `json:"Uploads_name"`
	OriginalName string    `json:"original_name"`
	UploadsPath  string    `json:"Uploads_path"`
	UploadsSize  int64     `json:"Uploads_size"`
	UploadsType  string    `json:"Uploads_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}
