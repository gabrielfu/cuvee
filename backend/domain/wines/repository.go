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

func (r *WineRepository) Get(ctx context.Context, id string) (WineDAO, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return WineDAO{}, err
	}
	var wine WineDAO
	res := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectId}})
	if err := res.Decode(&wine); err != nil {
		return WineDAO{}, err
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
	update := bson.D{{Key: "$set", Value: w}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *WineRepository) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err = r.collection.DeleteOne(context.Background(), filter)
	return err
}
