package repository

import (
	"Mongo/domain/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadsRepository interface {
	Create(Uploads *model.Uploads) error
	FindAll() ([]model.Uploads, error)
	FindByID(id string) (*model.Uploads, error)
	Delete(id string) error
}
type upRepository struct {
	collection *mongo.Collection
}

func NewUploadsRepository(db *mongo.Database) UploadsRepository {
	return &upRepository{
		collection: db.Collection("Uploads"),
	}
}
func (r *upRepository) Create(Uploads *model.Uploads) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Uploads.UploadedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, Uploads)
	if err != nil {
		return err
	}
	Uploads.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
func (r *upRepository) FindAll() ([]model.Uploads, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var Uploadss []model.Uploads
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &Uploadss); err != nil {
		return nil, err
	}
	return Uploadss, nil
}
func (r *upRepository) FindByID(id string) (*model.Uploads, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var Uploads model.Uploads
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&Uploads)
	if err != nil {
		return nil, err
	}
	return &Uploads, nil
}
func (r *upRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
