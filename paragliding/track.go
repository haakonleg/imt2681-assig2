package paragliding

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Retrieves a track from the database by its objectID (hex string)
func getTrackByID(db *Database, id string) (*Track, error) {
	objectID, err := objectid.FromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))
	tracks, err := db.findTracks(filter)
	if err != nil {
		return nil, err
	}

	if len(tracks) < 1 {
		return nil, errors.New("Track doesn't exist in database")
	}

	// Wtf, why is this even possible
	return &tracks[0], nil
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

// GET /api/track
// Returns an array of IDs of all tracks stored in the database
func getAllTracks(req *Request, db *Database) {
	// Get all tracks in database
	tracks, err := db.findTracks(nil)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	// Retrieve all the ids into a slice
	ids := make([]string, 0, len(tracks))
	for _, track := range tracks {
		ids = append(ids, track.ID.Hex())
	}

	req.SendJSON(&ids, http.StatusOK)
}

// GET /api/track/{id}
// Retrieves a track by the value of its ObjectID (hex encoded string)
func getTrack(req *Request, db *Database, id string) {
	track, err := getTrackByID(db, id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}
	req.SendJSON(track, http.StatusOK)
}

// GET /api/track/{id}/{field}
func getTrackField(req *Request, db *Database, id string, field string) {
	track, err := getTrackByID(db, id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	var resValue string
	switch field {
	case "pilot":
		resValue = track.Pilot
	case "glider":
		resValue = track.Glider
	case "glider_id":
		resValue = track.GliderID
	case "track_length":
		resValue = track.TrackLength
	case "H_date":
		resValue = track.HDate
	case "track_src_url":
		resValue = track.TrackSrcURL
	default:
		http.NotFound(req.w, req.r)
	}

	req.SendText(resValue)
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

	id, err := db.insertObject(newTrack, TRACKS)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)
}

// Routes the /track request to handlers
func handleTrackRequest(req *Request, db *Database, path string) {
	// GET/POST /api/track
	if match, _ := regexp.MatchString("^track[/]?$", path); match {
		switch req.r.Method {
		case "GET":
			getAllTracks(req, db)
		case "POST":
			registerTrack(req, db)
		}
		return
	}

	// GET /api/track/{id}
	if match := regexp.MustCompile("^track/([a-z0-9]{24})[/]?$").FindStringSubmatch(path); match != nil {
		getTrack(req, db, match[1])
		return
	}

	// GET track/{id}/{field}
	if match := regexp.MustCompile("^track/([a-z0-9]{24})/([a-zA-Z_]+)[/]?$").FindStringSubmatch(path); match != nil {
		getTrackField(req, db, match[1], match[2])
		return
	}

	http.NotFound(req.w, req.r)
}
