package paragliding

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

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

// GET /api/track
// Returns an array of IDs of all tracks stored in the database
func getAllTracks(req *Request, db *Database) {
	// Only get the id
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64("_id", 1)))}

	// Get all track IDs in database
	tracks, err := db.Find(TRACKS, nil, findopts)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	// Retrieve all the IDs into a slice
	ids := make([]string, 0, len(tracks))
	for _, track := range tracks {
		ids = append(ids, track.(Track).ID.Hex())
	}

	req.SendJSON(&ids, http.StatusOK)
}

// GET /api/track/{id}
// Retrieves a track by the value of its ObjectID (hex encoded string)
func getTrack(req *Request, db *Database, id string) {
	// Only get the requested track
	objectID, err := objectid.FromHex(id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))

	tracks, err := db.Find(TRACKS, filter, nil)
	if err != nil {
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
func getTrackField(req *Request, db *Database, id string, field string) {
	objectID, err := objectid.FromHex(id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	// Only get the requested field
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64(field, 1)))}

	tracks, err := db.Find(TRACKS, filter, findopts)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	track := tracks[0].(Track)
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
func registerTrack(req *Request, db *Database) {
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
	newTrack := createTrack(&igc, request.URL)
	id, err := db.InsertObject(TRACKS, &newTrack)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)

	// Invoke webhooks
	invokeWebhooks(db)
}
