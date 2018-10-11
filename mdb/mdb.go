package mdb

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// An enum of the database collections
type DatabaseCollection int

const (
	TRACKS DatabaseCollection = iota
	WEBHOOKS
)

// Stringer for databaseCollection type
func (dc DatabaseCollection) String() string {
	switch dc {
	case TRACKS:
		return "tracks"
	case WEBHOOKS:
		return "webhooks"
	}
	return ""
}

// Database contains the mongoDB database context, it also has helper methods for connecting to and querying the database
type Database struct {
	MongoURL string
	DBName   string

	client   *mongo.Client
	database *mongo.Database
}

// CreateConnection creates a connection to the mongoDB server
func (db *Database) CreateConnection() {
	client, err := mongo.Connect(context.Background(), db.MongoURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client
	db.database = db.client.Database(db.DBName)
	db.createTimestampIndex()
}

// InsertObject inserts an object into the specified collection in the database
func (db *Database) InsertObject(collection DatabaseCollection, object interface{}) (string, error) {
	col := db.database.Collection(collection.String())
	res, err := col.InsertOne(context.Background(), object)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), nil
}

// Find queries documents from the specified collection in the database
func (db *Database) Find(collection DatabaseCollection, filter interface{}, opts []findopt.Find, results interface{}) error {
	col := db.database.Collection(collection.String())
	cur, err := col.Find(context.Background(), filter, opts...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer cur.Close(context.Background())

	// Check that results is a slice and find its type
	switch resArr := results.(type) {
	case *[]*Track:
		elem := &Track{}
		for cur.Next(context.Background()) {
			if err := cur.Decode(elem); err != nil {
				return err
			}
			*resArr = append(*resArr, elem)
		}
	case *[]*Webhook:
		elem := &Webhook{}
		for cur.Next(context.Background()) {
			if err := cur.Decode(elem); err != nil {
				return err
			}
			*resArr = append(*resArr, elem)
		}
	default:
		log.Fatalf("This type is not supported: %s", reflect.TypeOf(resArr))
	}
	return nil
}

// Update updates documents in the specified collection in the database
func (db *Database) Update(collection DatabaseCollection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	col := db.database.Collection(collection.String())
	ur, err := col.UpdateMany(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ur, nil
}

// Creates a descending index on the timestamp field in tracks, to be able to support certain queries and better performance
func (db *Database) createTimestampIndex() {
	indexView := db.database.Collection(TRACKS.String()).Indexes()

	indexModel := mongo.IndexModel{
		Keys: bson.NewDocument(bson.EC.Int32("ts", -1))}

	_, err := indexView.CreateOne(context.Background(), indexModel, nil)
	if err != nil {
		log.Fatal(err)
	}
}
