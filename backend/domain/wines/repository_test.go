package wines_test

import (
	"context"
	"cuvee/domain/wines"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DATABASE   = "test_cuvee"
	COLLECTION = "wines"
)

var collection *mongo.Collection

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		uri := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}
		collection = client.Database(DATABASE).Collection(COLLECTION)
		return client.Ping(context.Background(), readpref.Primary())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func mustObjectId(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return objectId
}

func mockPurchase() *wines.MongoPurchase {
	return &wines.MongoPurchase{
		Quantity: 1,
		Price:    1000,
		Date:     "2020-01-01",
	}
}

func mockPurchase2() *wines.MongoPurchase {
	return &wines.MongoPurchase{
		Quantity: 2,
		Price:    2000,
		Date:     "2020-02-02",
	}
}

func mockWine() *wines.MongoWine {
	return &wines.MongoWine{
		Name:      "Chateau Margaux",
		Vintage:   "2015",
		Format:    "750ml",
		Country:   "France",
		Region:    "Bordeaux",
		Purchases: []wines.MongoPurchase{},
	}
}

type WinePurchaseIDs struct {
	WineID      string
	PurchaseIDs []string
}

func createWine(
	t *testing.T,
	collection *mongo.Collection,
	wine *wines.MongoWine,
	purchases ...*wines.MongoPurchase,
) WinePurchaseIDs {
	var out WinePurchaseIDs

	ctx := context.Background()
	wineRes, err := collection.InsertOne(ctx, wine)
	if err != nil {
		t.Fatalf("Insert wine failed: %s", err)
	}
	out.WineID = wineRes.InsertedID.(primitive.ObjectID).Hex()

	for _, purchase := range purchases {
		purchase = purchase.CopyWithNewID()
		filter := bson.M{"_id": mustObjectId(out.WineID)}
		update := bson.M{"$push": bson.M{"purchases": purchase}}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			t.Fatalf("Insert purchase failed: %s", err)
		}
		out.PurchaseIDs = append(out.PurchaseIDs, purchase.ID.Hex())
	}

	t.Cleanup(func() {
		deleteWine(t, collection, mustObjectId(out.WineID))
	})
	return out
}

func deleteWine(t *testing.T, collection *mongo.Collection, id primitive.ObjectID) {
	ctx := context.Background()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		t.Fatalf("Delete wine failed: %s", err)
	}
}

func TestCreateWine(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	wineId, err := repo.Create(ctx, wine)
	if err != nil {
		t.Fatalf("Create wine failed: %s", err)
	}

	assert.Nil(t, collection.FindOne(ctx, bson.M{"_id": mustObjectId(wineId)}).Err())
	deleteWine(t, collection, mustObjectId(wineId))
}

func TestListWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	createWine(t, collection, wine)

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List wines failed: %s", err)
	}

	assert.Equal(t, 1, len(list))
	assert.Equal(t, wine.Name, list[0].Name)
}

func TestGettWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	out := createWine(t, collection, wine)

	get, err := repo.Get(ctx, out.WineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestUpdateWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	out := createWine(t, collection, wine)

	wine.Name = "Chateau Latour"
	err := repo.Update(ctx, out.WineID, wine)
	if err != nil {
		t.Fatalf("Update wine failed: %s", err)
	}

	get, err := repo.Get(ctx, out.WineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestDeleteWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	out := createWine(t, collection, wine)

	err := repo.Delete(ctx, out.WineID)
	if err != nil {
		t.Fatalf("Delete wine failed: %s", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List wines failed: %s", err)
	}

	assert.Equal(t, 0, len(list))
}

func TestCreatePurchase(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	out := createWine(t, collection, wine)

	purchase := mockPurchase()
	purchaseId, err := repo.CreatePurchase(ctx, out.WineID, purchase)
	if err != nil {
		t.Fatalf("Create purchase failed: %s", err)
	}

	filter := bson.D{{Key: "_id", Value: mustObjectId(out.WineID)}, {Key: "purchases._id", Value: mustObjectId(purchaseId)}}
	assert.Nil(t, collection.FindOne(ctx, filter).Err())
}

func TestListPurchases(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	purchase := mockPurchase()
	out := createWine(t, collection, wine, purchase)

	list, err := repo.ListPurchases(ctx, out.WineID)
	if err != nil {
		t.Fatalf("List purchases failed: %s", err)
	}

	assert.Equal(t, 1, len(list))
	assert.Equal(t, out.PurchaseIDs[0], list[0].ID.Hex())
}

func TestGetPurchase(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	purchase := mockPurchase()
	out := createWine(t, collection, wine, purchase)

	get, err := repo.GetPurchase(ctx, out.PurchaseIDs[0])
	if err != nil {
		t.Fatalf("Get purchase failed: %s", err)
	}

	assert.Equal(t, out.PurchaseIDs[0], get.ID.Hex())
}

func TestUpdatePurchase(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	purchase := mockPurchase()
	out := createWine(t, collection, wine, purchase)

	newPurchase := mockPurchase2()
	err := repo.UpdatePurchase(ctx, out.PurchaseIDs[0], newPurchase)
	if err != nil {
		t.Fatalf("Update purchase failed: %s", err)
	}

	get, err := repo.GetPurchase(ctx, out.PurchaseIDs[0])
	if err != nil {
		t.Fatalf("Get purchase failed: %s", err)
	}

	assert.Equal(t, newPurchase.Quantity, get.Quantity)
}
