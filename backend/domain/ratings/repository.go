package ratings

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RegionDAO represents a region for a wine-vintage chart pair.
type RegionDAO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	WineID   string             `bson:"wine_id"`
	VCSymbol string             `bson:"vc_symbol"`
	Region   string             `bson:"region"`
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
	// Unique pair of wine_id and vc_symbol
	r.collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "wine_id", Value: 1}, {Key: "vc_symbol", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
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

func (r *RegionRepository) GetRegion(ctx context.Context, wineID, vcSymbol string) (*RegionDAO, error) {
	var result *RegionDAO
	res := r.collection.FindOne(ctx, bson.D{{Key: "wine_id", Value: wineID}, {Key: "vc_symbol", Value: vcSymbol}})
	if err := res.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RegionRepository) CreateRegion(ctx context.Context, region RegionDAO) error {
	_, err := r.collection.InsertOne(ctx, region)
	return err
}

func (r *RegionRepository) UpdateRegion(ctx context.Context, region RegionDAO) error {
	_, err := r.collection.ReplaceOne(ctx, bson.D{{Key: "wine_id", Value: region.WineID}, {Key: "vc_symbol", Value: region.VCSymbol}}, region)
	return err
}

func (r *RegionRepository) DeleteRegion(ctx context.Context, wineID, vcSymbol string) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{Key: "wine_id", Value: wineID}, {Key: "vc_symbol", Value: vcSymbol}})
	return err
}
