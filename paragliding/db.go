package paragliding

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// An enum of the database collections
type databaseCollection int

const (
	tracks databaseCollection = iota
	webhooks
)

// Stringer for databaseCollection type
func (dc databaseCollection) String() string {
	switch dc {
	case tracks:
		return "tracks"
	case webhooks:
		return "webhooks"
	default:
		return ""
	}
}

// Database contains the mongoDB database context, it also has helper methods for connecting to and querying the database
type Database struct {
	MongoURL string
	DBName   string

	client   *mongo.Client
	database *mongo.Database
	tracks   *mongo.Collection
	webhooks *mongo.Collection
}

// CreateConnection creates a connection to the mongoDB server
func (db *Database) createConnection() {
	client, err := mongo.Connect(context.Background(), db.MongoURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client
	db.database = db.client.Database(db.DBName)
	db.tracks = db.database.Collection(tracks.String())
	db.webhooks = db.database.Collection(webhooks.String())
	db.createTimestampIndex()
}

func (db *Database) insertObject(collection databaseCollection, object interface{}) (string, error) {
	col := db.database.Collection(collection.String())
	res, err := col.InsertOne(context.Background(), object)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), nil
}

func (db *Database) findTracks(filter interface{}, opts []findopt.Find) ([]Track, error) {
	cur, err := db.tracks.Find(context.Background(), filter, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	tracks := make([]Track, 0)
	for cur.Next(context.Background()) {
		var elem Track
		if err := cur.Decode(&elem); err != nil {
			fmt.Println(err)
			return nil, err
		}
		tracks = append(tracks, elem)
	}

	return tracks, nil
}

func (db *Database) findWebhooks(filter interface{}, opts []findopt.Find) ([]Webhook, error) {
	cur, err := db.webhooks.Find(context.Background(), filter, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	webhooks := make([]Webhook, 0)
	for cur.Next(context.Background()) {
		var elem Webhook
		if err := cur.Decode(&elem); err != nil {
			fmt.Println(err)
			return nil, err
		}
		webhooks = append(webhooks, elem)
	}

	return webhooks, nil
}

func (db *Database) updateWebhooks(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	ur, err := db.webhooks.UpdateMany(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ur, nil
}

// Creates a descending index on the timestamp field in tracks, to be able to support certain queries and better performance
func (db *Database) createTimestampIndex() {
	indexView := db.tracks.Indexes()

	indexModel := mongo.IndexModel{
		Keys: bson.NewDocument(bson.EC.Int32("ts", -1))}

	_, err := indexView.CreateOne(context.Background(), indexModel, nil)
	if err != nil {
		log.Fatal(err)
	}
}
