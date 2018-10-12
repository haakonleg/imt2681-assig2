package mdb

import (
	"strconv"
	"time"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Track is the model of IGC tracks stored in database
type Track struct {
	ID          objectid.ObjectID `bson:"_id" json:"-"`
	Ts          int64             `bson:"ts" json:"-"`
	HDate       string            `bson:"H_date" json:"H_Date"`
	Pilot       string            `bson:"pilot" json:"pilot"`
	Glider      string            `bson:"glider" json:"glider"`
	GliderID    string            `bson:"glider_id" json:"glider_id"`
	TrackLength string            `bson:"track_length" json:"track_length"`
	TrackSrcURL string            `bson:"track_src_url" json:"track_src_url"`
}

// Creates a new track object out of a parsed IGC track from goigc
func CreateTrack(igc *igc.Track, url string) Track {
	return Track{
		ID:          objectid.New(),
		Ts:          nowMilli(),
		HDate:       igc.Date.String(),
		Pilot:       igc.Pilot,
		Glider:      igc.GliderType,
		GliderID:    igc.GliderID,
		TrackLength: calTrackLen(igc.Points),
		TrackSrcURL: url}
}

func (t *Track) Field(field string) string {
	switch field {
	case "pilot":
		return t.Pilot
	case "glider":
		return t.Glider
	case "glider_id":
		return t.GliderID
	case "track_length":
		return t.TrackLength
	case "H_date":
		return t.HDate
	case "track_src_url":
		return t.TrackSrcURL
	default:
		return ""
	}
}

// Get current UNIX timestamp in miliseconds
func nowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Calculate track length, the sum of all distances between each point
func calTrackLen(points []igc.Point) string {
	d := 0.0
	for i := 0; i < len(points)-1; i++ {
		d += points[i].Distance(points[i+1])
	}
	return strconv.FormatFloat(d, 'f', 2, 64) + "km"
}
