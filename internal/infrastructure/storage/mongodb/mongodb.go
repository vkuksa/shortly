package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/link"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "shortly"
	collectionName = "links"
)

type Storage struct {
	collection *mongo.Collection
}

func NewStorage(mongoURI string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	collection := client.Database(databaseName).Collection(collectionName)
	err = createUUIDIndex(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("index creation failed: %w", err)
	}

	return &Storage{collection: collection}, nil
}

func createUUIDIndex(ctx context.Context, c *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "uuid", Value: 1}},
	}

	_, err := c.Indexes().CreateOne(ctx, indexModel, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetLink(ctx context.Context, uuid domain.UUID) (*domain.Link, error) {
	var l domain.Link
	err := s.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&l)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, link.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (s *Storage) StoreLink(ctx context.Context, link *domain.Link) error {
	_, err := s.collection.InsertOne(ctx, link)
	return err
}

func (s *Storage) IncHit(ctx context.Context, uuid domain.UUID) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"uuid": uuid},
		bson.M{"$inc": bson.M{"count": 1}},
	)
	if err != nil {
		return err
	}
	return nil
}
