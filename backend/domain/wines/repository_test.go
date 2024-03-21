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

func mockPurchase() *wines.MongoPurchase {
	return &wines.MongoPurchase{
		ID:       primitive.NewObjectID(),
		Quantity: 1,
		Price:    1000,
		Date:     "2020-01-01",
	}
}

func mockWine() *wines.MongoWine {
	return &wines.MongoWine{
		ID:        primitive.NewObjectID(),
		Name:      "Chateau Margaux",
		Vintage:   "2015",
		Format:    "750ml",
		Country:   "France",
		Region:    "Bordeaux",
		Purchases: []wines.MongoPurchase{*mockPurchase()},
	}
}

func createWine(t *testing.T, collection *mongo.Collection, wine *wines.MongoWine) {
	ctx := context.Background()
	_, err := collection.InsertOne(ctx, wine)
	if err != nil {
		t.Fatalf("Insert wine failed: %s", err)
	}

	t.Cleanup(func() {
		deleteWine(t, collection, wine.ID)
	})
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
	assert.Equal(t, wine.ID.Hex(), wineId)

	assert.Nil(t, collection.FindOne(ctx, bson.M{"_id": wine.ID}).Err())
	deleteWine(t, collection, wine.ID)
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
	createWine(t, collection, wine)

	get, err := repo.Get(ctx, wine.ID.Hex())
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestUpdateWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	createWine(t, collection, wine)

	wine.Name = "Chateau Latour"
	err := repo.Update(ctx, wine.ID.Hex(), wine)
	if err != nil {
		t.Fatalf("Update wine failed: %s", err)
	}

	get, err := repo.Get(ctx, wine.ID.Hex())
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestDeleteWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	createWine(t, collection, wine)

	err := repo.Delete(ctx, wine.ID.Hex())
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
	createWine(t, collection, wine)

	purchase := mockPurchase()
	purchaseId, err := repo.CreatePurchase(ctx, wine.ID.Hex(), purchase)
	if err != nil {
		t.Fatalf("Create purchase failed: %s", err)
	}
	assert.Equal(t, purchase.ID.Hex(), purchaseId)

	filter := bson.D{{Key: "_id", Value: wine.ID}, {Key: "purchases._id", Value: purchase.ID}}
	assert.Nil(t, collection.FindOne(ctx, filter).Err())
}

func TestListPurchases(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	createWine(t, collection, wine)

	list, err := repo.ListPurchases(ctx, wine.ID.Hex())
	if err != nil {
		t.Fatalf("List purchases failed: %s", err)
	}

	assert.Equal(t, 1, len(list))
	assert.Equal(t, wine.Purchases[0].ID.Hex(), list[0].ID.Hex())
}

func TestGetPurchase(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	createWine(t, collection, wine)

	get, err := repo.GetPurchase(ctx, wine.ID.Hex(), wine.Purchases[0].ID.Hex())
	if err != nil {
		t.Fatalf("Get purchase failed: %s", err)
	}
	log.Printf("%v\n", wine.Purchases[0])
	log.Printf("%v\n", get)

	assert.Equal(t, wine.Purchases[0].ID.Hex(), get.ID.Hex())
}
