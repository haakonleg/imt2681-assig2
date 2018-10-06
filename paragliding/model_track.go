package paragliding

import (
	"time"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Model of track stored in database
type track struct {
	ID          objectid.ObjectID `bson:"_id"`
	HDate       string            `json:"H_Date"`
	Pilot       string            `json:"pilot"`
	Glider      string            `json:"glider"`
	GliderID    string            `json:"glider_id"`
	TrackLength string            `json:"track_length"`
}

// Creates a new track object out of a parsed IGC track from goigc
func createTrack(igc *igc.Track) track {
	var newTrack track

	newTrack.ID = objectid.New()
	newTrack.HDate = igc.Date.String()
	newTrack.Pilot = igc.Pilot
	newTrack.Glider = igc.GliderType
	newTrack.GliderID = igc.GliderID
	newTrack.TrackLength = calTrackLen(&igc.Points).String()
	return newTrack
}

// Calculate track length, time of last point subtracted by time of first
func calTrackLen(points *[]igc.Point) time.Duration {
	arrLen := len(*points)

	firstTime := (*points)[0].Time
	lastTime := (*points)[arrLen-1].Time

	return lastTime.Sub(firstTime)
}
