package paragliding

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Db is the database context
type Db struct {
	client   *mongo.Client
	database *mongo.Database
}

// CreateConnection creates a connection to the mongoDB server
func (db *Db) CreateConnection(mongoURL string) error {
	client, err := mongo.Connect(context.Background(), mongoURL, nil)
	if err != nil {
		return err
	}
	db.client = client
	db.database = db.client.Database("imt2681-assig2")
	return nil
}
