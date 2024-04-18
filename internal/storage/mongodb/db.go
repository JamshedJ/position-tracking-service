package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/JamshedJ/position-tracking-service/internal/dto"
	"github.com/JamshedJ/position-tracking-service/protos/gen/pts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	GetNearbyPositions(ctx context.Context, latitude, longitude float64, p dto.PositionParams) ([]*pts.PositionResponse, error)
	UpdatePosition(ctx context.Context, req *pts.PositionRequest) error
}

func NewStorage(collection *mongo.Collection) Storage {
	return &storage{collection: collection}
}

type storage struct {
	collection *mongo.Collection
}

func New(uri, dbname, collName string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("error creating MongoDB client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	db := client.Database(dbname)          
	collection := db.Collection(collName)

	// Create base client
	initialClient := bson.M{
		"client_id": "1234",
		"latitude":  0,
		"longitude": 0,
	}

	_, err = collection.InsertOne(ctx, initialClient)
	if err != nil {
		log.Fatalf("error inserting initial client into MongoDB: %v", err)
	}

	log.Println("Initial client added successfully")

	return collection
}

func (s *storage) GetNearbyPositions(ctx context.Context, latitude, longitude float64, p dto.PositionParams) ([]*pts.PositionResponse, error) {
	filter := bson.M{
		"location": bson.M{
			"$geoWithin": bson.M{
				"$box": [][]float64{
					{p.MinLongitude, p.MinLatitude},
					{p.MaxLongitude, p.MaxLatitude},
				},
			},
		},
	}

	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching nearby positions from MongoDB: %v", err)

	}
	defer cur.Close(ctx)

	var results []*pts.PositionResponse
	for cur.Next(ctx) {
		var position pts.PositionResponse
		if err := cur.Decode(&position); err != nil {
			return nil, fmt.Errorf("error decoding nearby position: %v", err)

		}
		results = append(results, &position)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

func (s *storage) UpdatePosition(ctx context.Context, req *pts.PositionRequest) error {
	filter := bson.M{"client_id": req.ClientId}
	update := bson.M{
		"$set": bson.M{
			"location": bson.M{
				"type":        "Point",
				"coordinates": []float64{req.Longitude, req.Latitude},
			},
		},
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating position in MongoDB: %v", err)

	}
	return nil
}
