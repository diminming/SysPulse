package model

import (
	"context"
	"log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var Client driver.Client
var GraphDB driver.Database

func init() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		panic(err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "123456"),
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := client.Database(ctx, "insight")
	if err != nil {
		panic(err)
	}
	Client = client
	GraphDB = db
}

func CreateDocument(collection string, doc interface{}) {
	ctx := context.Background()
	coll, err := GraphDB.Collection(ctx, collection)
	if err != nil {
		panic(err)
	}
	meta, err := coll.CreateDocument(ctx, doc)
	if err != nil {
		panic(err)
	}
	log.Default().Printf("Created document with key '%s', revision '%s'\n", meta.Key, meta.Rev)
}
