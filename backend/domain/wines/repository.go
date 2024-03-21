package wines

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPurchase struct {
	ID       primitive.ObjectID `bson:"_id"`
	Quantity int                `bson:"quantity"`
	Price    float64            `bson:"price"`
	Date     string             `bson:"date"`
}

type MongoWine struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Vintage   string             `bson:"vintage"`
	Format    string             `bson:"format"`
	Country   string             `bson:"country"`
	Region    string             `bson:"region"`
	Purchases []MongoPurchase    `bson:"purchases"`
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

func (r *WineRepository) List(ctx context.Context) ([]*MongoWine, error) {
	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var wines []*MongoWine
	for cur.Next(ctx) {
		var w MongoWine
		if err := cur.Decode(&w); err != nil {
			return nil, err
		}
		wines = append(wines, &w)
	}
	return wines, nil
}

func (r *WineRepository) Get(ctx context.Context, id string) (*MongoWine, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var wine *MongoWine
	res := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectId}})
	if err := res.Decode(&wine); err != nil {
		return nil, err
	}
	return wine, nil
}

func (r *WineRepository) Create(ctx context.Context, w *MongoWine) (string, error) {
	_, err := r.collection.InsertOne(ctx, w)
	if err != nil {
		return "", err
	}
	return w.ID.Hex(), nil
}

func (r *WineRepository) Update(id string, w *MongoWine) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: w}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *WineRepository) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err = r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *WineRepository) ListPurchases(wineId string) ([]*MongoPurchase, error) {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return nil, err
	}
	var wine *MongoWine
	res := r.collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: objectId}})
	if err := res.Decode(&wine); err != nil {
		return nil, err
	}
	var out []*MongoPurchase
	for i := range wine.Purchases {
		out = append(out, &wine.Purchases[i])
	}
	return out, nil
}

func (r *WineRepository) GetPurchase(wineId, purchaseId string) (*MongoPurchase, error) {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return nil, err
	}
	var purchase *MongoPurchase
	filter := bson.D{{Key: "_id", Value: objectId}, {Key: "purchases._id", Value: purchaseId}}
	res := r.collection.FindOne(context.Background(), filter)
	if err := res.Decode(&purchase); err != nil {
		return nil, err
	}
	return purchase, nil
}

func (r *WineRepository) CreatePurchase(wineId string, p *MongoPurchase) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return "", err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "purchases", Value: p}}}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

func (r *WineRepository) UpdatePurchase(wineId, purchaseId string, p *MongoPurchase) error {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}, {Key: "purchases._id", Value: purchaseId}}
	update := bson.D{{Key: "$set", Value: p}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *WineRepository) DeletePurchase(wineId, purchaseId string) error {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "purchases._id", Value: purchaseId}}}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}
