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

// ValidateTrackID is the validation function used by the router to validate the track ID
// the ID is a hex encoded string consisting of 24 characters
func ValidateTrackID(variable string) (bool, interface{}) {
	validChars := "0123456789abcdef"
	if len(variable) != 24 {
		return false, nil
	}
	for _, ch := range variable {
		if !strings.ContainsRune(validChars, ch) {
			return false, nil
		}
	}
	return true, variable
}

// GET /api/track/{id}
// Retrieves a track by the value of its ObjectID (hex encoded string)
func (th *TrackHandler) GetTrack(req *router.Request) {
	id := req.Vars["id"].(string)

	// Only get the requested track
	objectID, err := objectid.FromHex(id)
	if err != nil {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))

	tracks := make([]*mdb.Track, 0)
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

// ValidateTrackField is the validator used by the router to validate a request for the field
// in a Track object
func ValidateTrackField(variable string) (bool, interface{}) {
	validFields := []string{
		"pilot", "glider",
		"glider_id", "track_length",
		"H_date", "track_src_url"}
	for _, field := range validFields {
		if variable == field {
			return true, variable
		}
	}
	return false, nil
}

// GET /api/track/{id}/{field}
func (th *TrackHandler) GetTrackField(req *router.Request) {
	id := req.Vars["id"].(string)
	field := req.Vars["field"].(string)

	objectID, _ := objectid.FromHex(id)

	// Only get the requested field
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64(field, 1)))}

	tracks := make([]*mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, filter, findopts, &tracks); err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("Invalid ID", http.StatusBadRequest)
		return
	}

	req.SendText(tracks[0].Field(field), http.StatusOK)
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
