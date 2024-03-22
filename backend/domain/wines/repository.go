package wines

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPurchase struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Quantity int                `bson:"quantity"`
	Price    float64            `bson:"price"`
	Date     string             `bson:"date"`
}

type MongoWine struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Vintage   string             `bson:"vintage"`
	Format    string             `bson:"format"`
	Country   string             `bson:"country"`
	Region    string             `bson:"region"`
	Purchases []MongoPurchase    `bson:"purchases"`
}

func (p *MongoPurchase) CopyWithNewID() *MongoPurchase {
	return &MongoPurchase{
		ID:       primitive.NewObjectID(),
		Quantity: p.Quantity,
		Price:    p.Price,
		Date:     p.Date,
	}
}

func (p *MongoPurchase) CopyWithoutID() *MongoPurchase {
	return &MongoPurchase{
		Quantity: p.Quantity,
		Price:    p.Price,
		Date:     p.Date,
	}
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
	res, err := r.collection.InsertOne(ctx, w)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *WineRepository) Update(ctx context.Context, id string, w *MongoWine) error {
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

func (r *WineRepository) ListPurchases(ctx context.Context, wineId string) ([]*MongoPurchase, error) {
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

func (r *WineRepository) GetPurchase(ctx context.Context, purchaseId string) (*MongoPurchase, error) {
	purchaseObjectId, err := primitive.ObjectIDFromHex(purchaseId)
	if err != nil {
		return nil, err
	}
	var wine *MongoWine
	filter := bson.D{{Key: "purchases._id", Value: purchaseObjectId}}
	res := r.collection.FindOne(context.Background(), filter)
	if err := res.Decode(&wine); err != nil {
		return nil, err
	}

	for i := range wine.Purchases {
		if wine.Purchases[i].ID == purchaseObjectId {
			return &wine.Purchases[i], nil
		}
	}
	return nil, fmt.Errorf("purchase %s not found", purchaseId)
}

func (r *WineRepository) CreatePurchase(ctx context.Context, wineId string, p *MongoPurchase) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(wineId)
	if err != nil {
		return "", err
	}
	p = p.CopyWithNewID()
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "purchases", Value: p}}}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "", err
	}
	// retrieve last purchase id
	return p.ID.Hex(), nil
}

func (r *WineRepository) UpdatePurchase(ctx context.Context, purchaseId string, p *MongoPurchase) error {
	old, err := r.GetPurchase(ctx, purchaseId)
	if err != nil {
		return err
	}

	purchaseObjectId, err := primitive.ObjectIDFromHex(purchaseId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "purchases._id", Value: purchaseObjectId}}
	new := p.CopyWithoutID()
	new.ID = old.ID
	update := bson.D{{Key: "$set", Value: bson.M{"purchases.$": new}}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *WineRepository) DeletePurchase(ctx context.Context, purchaseId string) error {
	purchaseObjectId, err := primitive.ObjectIDFromHex(purchaseId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "purchases._id", Value: purchaseObjectId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "purchases", Value: bson.M{"_id": purchaseObjectId}}}}}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}
