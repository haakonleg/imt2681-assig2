package main

import (
	"os"
	"strings"

	"github.com/haakonleg/imt2681-assig2/paragliding"
)

// API config
const dbUser = "testuser1"
const dbPass = "testpass1"
const mongoURL = "mongodb://<dbuser>:<dbpassword>@ds223063.mlab.com:23063/imt2681-assig2"
const dbName = "imt2681-assig2"

const defaultPort = "8080"

// Main starts the paragliding server by supplying the configuration options to the App object
func main() {
	url := strings.Replace(mongoURL, "<dbuser>", dbUser, 1)
	url = strings.Replace(url, "<dbpassword>", dbPass, 1)

	// Get port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	// Configure and start the API
	app := paragliding.App{
		MongoURL:    url,
		ListenPort:  port,
		DBName:      dbName,
		TickerLimit: 5}
	app.StartServer()
}
