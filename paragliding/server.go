package paragliding

import (
	"fmt"
	"log"
	"net/http"

	"github.com/haakonleg/imt2681-assig2/admin"
	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/haakonleg/imt2681-assig2/ticker"
	"github.com/haakonleg/imt2681-assig2/track"
	"github.com/haakonleg/imt2681-assig2/webhook"
)

// App must be instantiated with the url to the mongodb database, database name and the port for the API to listen on
type App struct {
	MongoURL    string
	DBName      string
	ListenPort  string
	TickerLimit int64

	db             *mdb.Database
	infoHandler    *ApiInfoHandler
	trackHandler   *track.TrackHandler
	tickerHandler  *ticker.TickerHandler
	webhookHandler *webhook.WebhookHandler
	adminHandler   *admin.AdminHandler
}

func (app *App) configureRoutes(r *router.Router) {
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

	// Webhook routes
	r.Handle("POST", "/paragliding/api/webhook/new_track", app.webhookHandler.PostWebhook)
	r.Handle("GET", "/paragliding/api/webhook/new_track/{id}", app.webhookHandler.GetWebhook)
	r.Handle("DELETE", "/paragliding/api/webhook/new_track/{id}", app.webhookHandler.DeleteWebhook)

	// Admin routes
	r.Handle("GET", "/admin/api/tracks_count", app.adminHandler.GetTrackCount)
	r.Handle("DELETE", "/admin/api/tracks", app.adminHandler.DeleteAllTracks)
}

func (app *App) configureValidators(r *router.Router) {
	r.Validate("id", track.ValidateTrackID)
	r.Validate("field", track.ValidateTrackField)
	r.Validate("timestamp", ticker.ValidateTimestamp)
}

// StartServer starts listening and serving the API server
func (app *App) StartServer() {
	// Try connect to mongoDB
	app.db = &mdb.Database{MongoURL: app.MongoURL, DBName: app.DBName}
	app.db.CreateConnection()
	fmt.Println("Connected to mongoDB")

	// Create handlers
	app.infoHandler = NewInfoHandler()
	app.trackHandler = track.NewTrackHandler(app.db)
	app.tickerHandler = ticker.NewTickerHandler(app.TickerLimit, app.db)
	app.webhookHandler = webhook.NewWebhookHandler(app.db)
	app.adminHandler = admin.NewAdminHandler(app.db)

	// Registers a callback so that when a new track is registered, the webhook handler will check
	// if any webhooks should be triggered
	app.trackHandler.SetTrackRegisterCallback(app.webhookHandler.CheckInvokeWebhooks)

	// Instantiate router, and configure the handlers and paths
	r := router.NewRouter()
	app.configureRoutes(r)
	app.configureValidators(r)

	// Start listen
	fmt.Printf("Server listening on port %s\n", app.ListenPort)
	if err := http.ListenAndServe(":"+app.ListenPort, r); err != nil {
		log.Fatal(err.Error())
	}
}
