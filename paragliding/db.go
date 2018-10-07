package paragliding

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// Database contains the mongoDB database context, it also has helper methods for connecting to and querying the database
type Database struct {
	MongoURL string
	DBName   string

	client   *mongo.Client
	database *mongo.Database
	tracks   *mongo.Collection
}

// CreateConnection creates a connection to the mongoDB server
func (db *Database) createConnection() {
	client, err := mongo.Connect(context.Background(), db.MongoURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client
	db.database = db.client.Database(db.DBName)
	db.tracks = db.database.Collection("tracks")
	db.createTimestampIndex()
}

func (db *Database) insertTrack(track *Track) (string, error) {
	res, err := db.tracks.InsertOne(context.Background(), track)
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
