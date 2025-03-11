package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vkuksa/shortly/internal/link"
)

const (
	collectionName = "links"
)

type Storage struct {
	collection *mongo.Collection
}

func NewStorage(db *mongo.Database) (*Storage, error) {
	coll := db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := createIndexes(ctx, coll); err != nil {
		return nil, fmt.Errorf("createIndexes: %w", err)
	}

	return &Storage{collection: coll}, nil
}

func createIndexes(ctx context.Context, c *mongo.Collection) error {
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "shortened", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := c.Indexes().CreateOne(ctx, idx)
	return err
}

func (s *Storage) Get(ctx context.Context, shortened string) (*link.ShortenedLink, error) {
	filter := bson.M{"shortened": shortened}

	var doc link.ShortenedLink
	err := s.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, link.ErrNotFound
		}
		return nil, fmt.Errorf("mongo findOne: %w", err)
	}

	return &doc, nil
}

// Store either inserts or replaces a link document. If a doc with the same
// `_id` already exists, it will be replaced; otherwise, it will be inserted.
func (s *Storage) Store(ctx context.Context, l *link.ShortenedLink) error {
	// We use ReplaceOne with upsert = true, keyed by _id
	filter := bson.M{"_id": l.Id}
	opts := options.Replace().SetUpsert(true)
	_, err := s.collection.ReplaceOne(ctx, filter, l, opts)
	if err != nil {
		return err
	}

	return nil
}

// AddHit increments the hits for a given link (by UUID).
func (s *Storage) AddHit(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"hits": 1}}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return link.ErrNotFound
	}
	return nil
}

// UpdateExpiration modifies the expiration time of a given link (by UUID).
func (s *Storage) UpdateExpiration(ctx context.Context, id string, expiresAt time.Time) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"expiresAt": expiresAt}}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return link.ErrNotFound
	}
	return nil
}
