package mongodb

import (
	"context"
	"time"

	"80.GO/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderMongoDb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewOrder() (*OrderMongoDb, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("Orders").Collection("order")
	return &OrderMongoDb{client: client, collection: collection}, nil
}

func (o *OrderMongoDb) CreateOrderMongoDb(order models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := o.collection.InsertOne(ctx, order)
	return err
}

func (o *OrderMongoDb) GetOrderMongoDb(id primitive.ObjectID) (*models.Order, error) {
	var order models.Order

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := o.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
