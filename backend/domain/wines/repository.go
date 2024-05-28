package wines

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PurchaseDAO struct {
	Quantity int     `bson:"quantity"`
	Price    float64 `bson:"price"`
	Date     string  `bson:"date"`
}

type WineDAO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Vintage   string             `bson:"vintage"`
	Format    string             `bson:"format"`
	Country   string             `bson:"country"`
	Region    string             `bson:"region"`
	Purchases []PurchaseDAO      `bson:"purchases"`
	ImageUrl  string             `bson:"image_url"`
}

// WineRepository implements WineRepository
type WineRepository struct {
	collection *mongo.Collection
}

func NewWineRepository(collection *mongo.Collection) *WineRepository {
	return &WineRepository{
		collection: collection,
	}
}

func (r *WineRepository) List(ctx context.Context) ([]WineDAO, error) {
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var wines []WineDAO
	if err = cur.All(ctx, &wines); err != nil {
		return nil, err
	}
	return wines, nil
}

func (r *WineRepository) Get(ctx context.Context, id string) (*WineDAO, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var wine *WineDAO
	res := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectId}})
	if err := res.Decode(&wine); err != nil {
		return nil, err
	}
	return wine, nil
}

func (r *WineRepository) Create(ctx context.Context, w *WineDAO) (string, error) {
	res, err := r.collection.InsertOne(ctx, w)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *WineRepository) Update(ctx context.Context, id string, w *WineDAO) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}

	var update bson.D
	if w.Name != "" {
		update = append(update, bson.E{Key: "name", Value: w.Name})
	}
	if w.Vintage != "" {
		update = append(update, bson.E{Key: "vintage", Value: w.Vintage})
	}
	if w.Format != "" {
		update = append(update, bson.E{Key: "format", Value: w.Format})
	}
	if w.Country != "" {
		update = append(update, bson.E{Key: "country", Value: w.Country})
	}
	if w.Region != "" {
		update = append(update, bson.E{Key: "region", Value: w.Region})
	}
	if w.ImageUrl != "" {
		update = append(update, bson.E{Key: "image_url", Value: w.ImageUrl})
	}
	if w.Purchases != nil {
		update = append(update, bson.E{Key: "purchases", Value: w.Purchases})
	}
	if len(update) == 0 {
		return nil
	}

	_, err = r.collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	return err
}

func (r *WineRepository) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
