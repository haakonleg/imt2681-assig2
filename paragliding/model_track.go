package paragliding

import (
	"time"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Track is the model of IGC tracks stored in database
type Track struct {
	ID          objectid.ObjectID `bson:"_id"`
	HDate       string            `bson:"H_date" json:"H_Date"`
	Pilot       string            `bson:"pilot" json:"pilot"`
	Glider      string            `bson:"glider" json:"glider"`
	GliderID    string            `bson:"glider_id" json:"glider_id"`
	TrackLength string            `bson:"track_length" json:"track_length"`
}

// Creates a new track object out of a parsed IGC track from goigc
func createTrack(igc *igc.Track) Track {
	return Track{
		ID:          objectid.New(),
		HDate:       igc.Date.String(),
		Pilot:       igc.Pilot,
		Glider:      igc.GliderType,
		GliderID:    igc.GliderID,
		TrackLength: calTrackLen(&igc.Points).String()}
}

// Calculate track length, time of last point subtracted by time of first
func calTrackLen(points *[]igc.Point) time.Duration {
	arrLen := len(*points)

	firstTime := (*points)[0].Time
	lastTime := (*points)[arrLen-1].Time

	return lastTime.Sub(firstTime)
}
