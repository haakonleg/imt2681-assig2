package paragliding

import (
	"fmt"
	"log"
	"net/http"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/haakonleg/imt2681-assig2/ticker"
	"github.com/haakonleg/imt2681-assig2/track"
)

// App must be instantiated with the url to the mongodb database, database name and the port for the API to listen on
type App struct {
	MongoURL    string
	DBName      string
	ListenPort  string
	TickerLimit int64

	db            *mdb.Database
	infoHandler   *ApiInfoHandler
	trackHandler  *track.TrackHandler
	tickerHandler *ticker.TickerHandler
}

// StartServer starts listening and serving the API server
func (app *App) StartServer() {

	// Try connect to mongoDB
	app.db = &mdb.Database{MongoURL: app.MongoURL, DBName: app.DBName}
	app.db.createConnection()
	fmt.Println("Connected to mongoDB")

	// Create handlers
	app.infoHandler = NewInfoHandler()
	app.trackHandler = NewTrackHandler(app.db)
	app.tickerHandler = NewTickerHandler(app.TickerLimit, app.db)

	// Instantiate router, and configure the paths
	r := router.NewRouter()

	// Redirect to /paragliding/api
	r.Handle("GET", "/paragliding", func(req *router.Request) {
		req.Redirect("/paragliding/api")
	})

	// Track routes
	r.Handle("GET", "/paragliding/api", app.infoHandler.getAPIInfo)
	r.Handle("POST", "/paragliding/api/track", app.trackHandler.PostTrack)
	r.Handle("GET", "/paragliding/api/track", app.trackHandler.GetAllTracks)
	r.Handle("GET", "/paragliding/api/track/{id}", app.trackHandler.GetTrack)
	r.Handle("GET", "/paragliding/api/track/{id}/{field}", app.trackHandler.GetTrackField)

	// Ticker routes
	r.Handle("GET", "/paragliding/api/ticker/latest", app.tickerHandler.GetLatestTimestamp)
	r.Handle("GET", "/paragliding/api/ticker", app.tickerHandler.GetTicker)
	r.Handle("GET", "/paragliding/api/ticker/{timestamp}", app.tickerHandler.GetTicker)

	// Start listen
	fmt.Printf("Server listening on port %s\n", app.ListenPort)
	if err := http.ListenAndServe(":"+app.ListenPort, &r); err != nil {
		log.Fatal(err.Error())
	}
}
