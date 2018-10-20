package main

import (
	"log"
	"os"

	"github.com/haakonleg/imt2681-assig2/paragliding"
)

// API config
const (
	dbName      = "imt2681-assig2"
	defaultPort = "8080"
)

// Main starts the paragliding server by supplying the configuration options to the App object
func main() {
	// Get port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	mongoURL := os.Getenv("PARAGLIDING_MONGO")
	if len(mongoURL) == 0 {
		log.Fatal("PARAGLIDING_MONGO environment variable is not set (put mongodb url in here)")
	}

	// Configure and start the API
	app := paragliding.App{
		MongoURL:    mongoURL,
		ListenPort:  port,
		DBName:      dbName,
		TickerLimit: 5}
	app.StartServer()
}
