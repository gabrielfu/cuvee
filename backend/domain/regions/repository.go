package regions

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RegionDAO represents a region for a wine-vintage chart pair.
type RegionDAO struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	WineID string             `bson:"wine_id"`
	Symbol string             `bson:"symbol"`
	Region string             `bson:"region"`
}

type RegionRepository struct {
	collection *mongo.Collection
}

func NewRegionRepository(ctx context.Context, collection *mongo.Collection) *RegionRepository {
	r := &RegionRepository{
		collection: collection,
	}
	r.createIndexes(ctx)
	return r
}

func (r *RegionRepository) createIndexes(ctx context.Context) {
	// Unique pair of wine_id and symbol
	r.collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "wine_id", Value: 1}, {Key: "symbol", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
}

func makeFilter(wineID, symbol string) bson.D {
	return bson.D{{Key: "wine_id", Value: wineID}, {Key: "symbol", Value: symbol}}
}

func (r *RegionRepository) ListRegions(ctx context.Context, wineID string) ([]RegionDAO, error) {
	cur, err := r.collection.Find(ctx, bson.D{{Key: "wine_id", Value: wineID}})
	if err != nil {
		return nil, err
	}
	var result []RegionDAO
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RegionRepository) GetRegion(ctx context.Context, wineID, symbol string) (*RegionDAO, error) {
	var result *RegionDAO
	filter := makeFilter(wineID, symbol)
	res := r.collection.FindOne(ctx, filter)
	if err := res.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RegionRepository) CreateRegion(ctx context.Context, region RegionDAO) error {
	_, err := r.collection.InsertOne(ctx, region)
	return err
}

func (r *RegionRepository) UpdateRegion(ctx context.Context, wineID, symbol, region string) error {
	filter := makeFilter(wineID, symbol)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "region", Value: region}}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *RegionRepository) DeleteRegion(ctx context.Context, wineID, symbol string) error {
	filter := makeFilter(wineID, symbol)
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
