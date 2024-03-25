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

func mockPurchase() wines.PurchaseDAO {
	return wines.PurchaseDAO{
		Quantity: 1,
		Price:    1000,
		Date:     "2020-01-01",
	}
}

func mockPurchase2() wines.PurchaseDAO {
	return wines.PurchaseDAO{
		Quantity: 2,
		Price:    2000,
		Date:     "2020-02-02",
	}
}

func mockPurchase3() wines.PurchaseDAO {
	return wines.PurchaseDAO{
		Quantity: 3,
		Price:    3000,
		Date:     "2020-03-03",
	}
}

func mockWine() wines.WineDAO {
	return wines.WineDAO{
		Name:      "Chateau Margaux",
		Vintage:   "2015",
		Format:    "750ml",
		Country:   "France",
		Region:    "Bordeaux",
		Purchases: []wines.PurchaseDAO{},
	}
}

func mockWineWithPurchases() wines.WineDAO {
	return wines.WineDAO{
		Name:    "Chateau Margaux",
		Vintage: "2015",
		Format:  "750ml",
		Country: "France",
		Region:  "Bordeaux",
		Purchases: []wines.PurchaseDAO{
			mockPurchase(),
			mockPurchase2(),
		},
	}
}

type WinePurchaseIDs struct {
	WineID      string
	PurchaseIDs []string
}

// returns the WineDAO id as a string
func createWine(
	t *testing.T,
	collection *mongo.Collection,
	wine wines.WineDAO,
	purchases ...wines.PurchaseDAO,
) string {
	ctx := context.Background()
	wineRes, err := collection.InsertOne(ctx, wine)
	if err != nil {
		t.Fatalf("Insert wine failed: %s", err)
	}
	wineID := wineRes.InsertedID.(primitive.ObjectID).Hex()

	for _, purchase := range purchases {
		filter := bson.M{"_id": mustObjectId(wineID)}
		update := bson.M{"$push": bson.M{"purchases": purchase}}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			t.Fatalf("Insert purchase failed: %s", err)
		}
	}

	t.Cleanup(func() {
		deleteWine(t, collection, mustObjectId(wineID))
	})
	return wineID
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
	wineId, err := repo.Create(ctx, &wine)
	if err != nil {
		t.Fatalf("Create wine failed: %s", err)
	}

	assert.Nil(t, collection.FindOne(ctx, bson.M{"_id": mustObjectId(wineId)}).Err())
	deleteWine(t, collection, mustObjectId(wineId))
}

func TestCreateWineWithPurchases(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWineWithPurchases()
	wineId, err := repo.Create(ctx, &wine)
	if err != nil {
		t.Fatalf("Create wine failed: %s", err)
	}

	assert.Nil(t, collection.FindOne(ctx, bson.M{"_id": mustObjectId(wineId)}).Err())
	assert.Equal(t, 2, len(wine.Purchases))
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

func TestGetWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	wineID := createWine(t, collection, wine)

	get, err := repo.Get(ctx, wineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestUpdateWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	wineID := createWine(t, collection, wine)

	wine.Name = "Chateau Latour"
	err := repo.Update(ctx, wineID, &wine)
	if err != nil {
		t.Fatalf("Update wine failed: %s", err)
	}

	get, err := repo.Get(ctx, wineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, wine.Name, get.Name)
}

func TestDeleteWines(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	wineID := createWine(t, collection, wine)

	err := repo.Delete(ctx, wineID)
	if err != nil {
		t.Fatalf("Delete wine failed: %s", err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List wines failed: %s", err)
	}

	assert.Equal(t, 0, len(list))
}

func TestAddPurchaseToWine(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWine()
	wineID := createWine(t, collection, wine)

	purchase := mockPurchase()
	wine.Purchases = append(wine.Purchases, purchase)

	err := repo.Update(ctx, wineID, &wine)
	if err != nil {
		t.Fatalf("Create purchase failed: %s", err)
	}

	get, err := repo.Get(ctx, wineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, 1, len(get.Purchases))
	assert.Equal(t, purchase.Quantity, get.Purchases[0].Quantity)
}

func TestModifyPurchasesInWine(t *testing.T) {
	repo := wines.NewWineRepository(collection)
	ctx := context.Background()

	wine := mockWineWithPurchases()
	wineID := createWine(t, collection, wine)

	purchase := mockPurchase3()
	wine.Purchases = []wines.PurchaseDAO{purchase}

	err := repo.Update(ctx, wineID, &wine)
	if err != nil {
		t.Fatalf("Create purchase failed: %s", err)
	}

	get, err := repo.Get(ctx, wineID)
	if err != nil {
		t.Fatalf("Get wine failed: %s", err)
	}

	assert.Equal(t, 1, len(get.Purchases))
	assert.Equal(t, purchase.Quantity, get.Purchases[0].Quantity)
}
