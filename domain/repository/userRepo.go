package repository

import (
	"Mongo/domain/config"
	"Mongo/domain/model"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionUsers = "users"

type userRepoStruct struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) model.UserRepository {
	return &userRepoStruct{client}
}

func (r *userRepoStruct) getCollection() *mongo.Collection {
	return r.client.Database("alumni_management_db").Collection(CollectionUsers)
}

func (r *userRepoStruct) FindByID(id primitive.ObjectID) (*model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.getCollection()
	user := new(model.Users)

	filter := bson.M{"_id": id}

	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepoStruct) FindByEmail(email string) (*model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.getCollection()
	user := new(model.Users)

	filter := bson.M{"email": email}

	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepoStruct) FindAll() ([]model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.getCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.Users
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepoStruct) Create(user *model.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.getCollection()

	user.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

func (r *userRepoStruct) Update(user *model.Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.getCollection()

	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"email":    user.Email,
			"username": user.Username,
			"password": user.Password,
			"role":     user.Role,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *userRepoStruct) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.getCollection()
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *userRepoStruct) Count(search string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := r.getCollection()

	filter := bson.M{}
	if search != "" {
		searchPattern := primitive.Regex{Pattern: search, Options: "i"}
		filter = bson.M{
			"$or": []bson.M{
				{"username": searchPattern},
				{"email": searchPattern},
			},
		}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func GetUsersRepo(search, sortBy, order string, limit, offset int) ([]model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.DB.Database("alumni_management_db").Collection(CollectionUsers)

	filter := bson.M{}
	if search != "" {
		searchPattern := primitive.Regex{Pattern: search, Options: "i"}
		filter = bson.M{
			"$or": []bson.M{
				{"username": searchPattern},
				{"email": searchPattern},
			},
		}
	}

	sortDirection := 1
	if order == "desc" {
		sortDirection = -1
	}
	sort := bson.D{{Key: sortBy, Value: sortDirection}}

	findOptions := options.Find().
		SetSort(sort).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.Users
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func CountUsersRepo(search string) (int, error) {

	repo := NewUserRepository(config.DB)

	return repo.Count(search)
}
