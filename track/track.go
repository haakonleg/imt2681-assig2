package track

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type TrackHandler struct {
	db *mdb.Database
}

func NewTrackHandler(db *mdb.Database) *TrackHandler {
	return &TrackHandler{
		db: db}
}

// GET /api/track
// Returns an array of IDs of all tracks stored in the database
func (th *TrackHandler) GetAllTracks(req *router.Request) {
	// Only get the id
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64("_id", 1)))}

	// Get all track IDs in database
	tracks := make([]*mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, nil, findopts, &tracks); err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	// Retrieve all the IDs into a slice
	ids := make([]string, 0, len(tracks))
	for _, track := range tracks {
		ids = append(ids, track.ID.Hex())
	}

	req.SendJSON(&ids, http.StatusOK)
}

// GET /api/track/{id}
// Retrieves a track by the value of its ObjectID (hex encoded string)
func (th *TrackHandler) GetTrack(req *router.Request) {
	id := req.Vars[0]

	// Only get the requested track
	objectID, err := objectid.FromHex(id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))

	tracks := make([]mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, filter, nil, &tracks); err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	req.SendJSON(&tracks[0], http.StatusOK)
}

// GET /api/track/{id}/{field}
func (th *TrackHandler) GetTrackField(req *router.Request) {
	id := req.Vars[0]
	field := req.Vars[1]

	objectID, err := objectid.FromHex(id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	// Only get the requested field
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64(field, 1)))}

	tracks := make([]mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, filter, findopts, &tracks); err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	track := tracks[0]
	switch field {
	case "pilot":
		req.SendText(track.Pilot)
	case "glider":
		req.SendText(track.Glider)
	case "glider_id":
		req.SendText(track.GliderID)
	case "track_length":
		req.SendText(track.TrackLength)
	case "H_date":
		req.SendText(track.HDate)
	case "track_src_url":
		req.SendText(track.TrackSrcURL)
	}
}

// POST /api/track
// Register/upload a track
func (th *TrackHandler) PostTrack(req *router.Request) {
	var request struct {
		URL string `json:"url"`
	}

	// Get the JSON post request
	if err := req.ParseJSONRequest(&request); err != nil {
		req.SendError("Error parsing JSON request", http.StatusBadRequest)
		return
	}

	// Check that the supplied link is valid
	if valid := ensureIGCLink(request.URL); !valid {
		req.SendError("This is not a valid IGC link", http.StatusBadRequest)
		return
	}

	// Parse the IGC file
	igc, err := igc.ParseLocation(request.URL)
	if err != nil {
		fmt.Println(err)
		req.SendError("Error parsing IGC track", http.StatusBadRequest)
		return
	}

	// Send response containing the ID to the inserted track
	newTrack := mdb.CreateTrack(&igc, request.URL)
	id, err := th.db.InsertObject(mdb.TRACKS, &newTrack)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)

	// Invoke webhooks
	//checkInvokeWebhooks(th.db)
}

// Ensures that a link points to an IGC resource (but just that it is a valid URL and has an igc extension)
func ensureIGCLink(link string) bool {
	if _, err := url.ParseRequestURI(link); err != nil {
		return false
	}

	ext := strings.ToLower(path.Ext(link))
	if ext != ".igc" {
		return false
	}
	return true
}
