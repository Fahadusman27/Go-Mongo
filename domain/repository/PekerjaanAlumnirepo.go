package repository

import (
	"context"
	"errors"
	"tugas/domain/config"
	"tugas/domain/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionPekerjaan = "pekerjaan_alumni"

func getCollectionPekerjaan() *mongo.Collection {
	return config.DB.Database("mahasiswa").Collection(CollectionPekerjaan)
}

func CheckpekerjaanAlumniByID(id string) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid job ID format")
	}

	filter := bson.M{"_id": objID}

	err = collection.FindOne(ctx, filter).Decode(pekerjaan)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func CreatepekerjaanAlumni(pekerjaan *model.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	pekerjaan.CreatedAt = time.Now()
	pekerjaan.UpdatedAt = time.Now()
	
	result, err := collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return err
	}

    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        pekerjaan.ID = oid
    } else {
        return errors.New("failed to get inserted ID")
    }

	return nil
}

func UpdatepekerjaanAlumni(NimAlumni string, pekerjaan *model.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	filter := bson.M{"nim_alumni": NimAlumni}

	update := bson.M{
		"$set": bson.M{
			"status_kerja": pekerjaan.StatusKerja,
			"jenis_industri": pekerjaan.JenisIndustri,
			"pekerjaan": pekerjaan.Pekerjaan,
			"jabatan": pekerjaan.Jabatan,
			"gaji": pekerjaan.Gaji,
			"lama_bekerja": pekerjaan.LamaBekerja,
			"updated_at": time.Now(),
		},
	}
    
    opts := options.Update().SetUpsert(false) 

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return errors.New("no document found or updated for the given nim")
	}

	return nil
}

func GetAllpekerjaanAlumni() ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

    filter := bson.M{"is_deleted": bson.M{"$exists": false}}
	
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}
	return pekerjaanList, nil
}

func SoftDeleteBynim(NimAlumni string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	filter := bson.M{"nim_alumni": NimAlumni, "is_deleted": bson.M{"$exists": false}}
	
	update := bson.M{"$set": bson.M{"is_deleted": time.Now()}}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.ModifiedCount == 0 {
        return errors.New("no active job record found for the given nim")
    }
	
	return nil
}


func GetAllTrash(nimAlumni string) ([]*model.Trash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	filter := bson.M{"is_deleted": bson.M{"$exists": true}}

	if nimAlumni != "" {
		filter["nim_alumni"] = nimAlumni
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trashes []*model.Trash
	if err = cursor.All(ctx, &trashes); err != nil {
		return nil, err
	}
	
	return trashes, nil
}

func RestoreTrashBynim(NimAlumni string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	filter := bson.M{"nim_alumni": NimAlumni, "is_deleted": bson.M{"$exists": true}}
	
	update := bson.M{"$unset": bson.M{"is_deleted": ""}}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
        return errors.New("no deleted job record found for the given nim")
    }

	return nil
}

func DeletePekerjaanByid(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := getCollectionPekerjaan()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid job ID format")
	}

	filter := bson.M{"_id": objID, "is_deleted": bson.M{"$exists": true}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no deleted document found with that ID")
	}

	return nil
}