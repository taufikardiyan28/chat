package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Conn struct {
		Pool *mongo.Client
	}
)

var Pool *Conn
var mongoClient *mongo.Client

func (c *Conn) Connect(uri string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	c.Pool = client
	mongoClient = client
	return err
}

func (c *Conn) GetPool() *mongo.Client {
	return mongoClient
}
