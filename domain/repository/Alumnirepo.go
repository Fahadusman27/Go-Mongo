package repository

import (
	. "Mongo/domain/config"
	"Mongo/domain/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAlumniCollection() *mongo.Collection {
	return DB.Database("alumni_management_db").Collection("alumni")
}

func CheckAlumniByNim(nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	ctx := context.TODO()
	collection := getAlumniCollection()

	filter := bson.M{"nim": nim}

	err := collection.FindOne(ctx, filter).Decode(alumni)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func CreateAlumni(alumni *model.Alumni) error {
	ctx := context.TODO()
	collection := getAlumniCollection()

	_, err := collection.InsertOne(ctx, alumni)
	return err
}

func UpdateAlumni(nim string, alumni *model.Alumni) error {
	ctx := context.TODO()
	collection := getAlumniCollection()

	filter := bson.M{"nim": nim}

	update := bson.M{
		"$set": alumni,
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func DeleteAlumni(nim string) error {
	ctx := context.TODO()
	collection := getAlumniCollection()

	filter := bson.M{"nim": nim}

	_, err := collection.DeleteOne(ctx, filter)
	return err
}

func GetAllAlumni() ([]model.Alumni, error) {
	ctx := context.TODO()
	collection := getAlumniCollection()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}

	return alumniList, nil
}
