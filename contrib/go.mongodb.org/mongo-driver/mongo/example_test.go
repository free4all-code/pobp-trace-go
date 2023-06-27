

package mongo_test

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongotrace "git.proto.group/protoobp/pobp-trace-go/contrib/go.mongodb.org/mongo-driver/mongo"
)

func Example() {
	// connect to MongoDB
	opts := options.Client()
	opts.Monitor = mongotrace.NewMonitor()
	opts.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("example")
	inventory := db.Collection("inventory")

	inventory.InsertOne(context.Background(), bson.D{
		{Key: "item", Value: "canvas"},
		{Key: "qty", Value: 100},
		{Key: "tags", Value: bson.A{"cotton"}},
		{Key: "size", Value: bson.D{
			{Key: "h", Value: 28},
			{Key: "w", Value: 35.5},
			{Key: "uom", Value: "cm"},
		}},
	})
}
