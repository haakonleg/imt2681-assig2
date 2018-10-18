package mdb

import (
	"github.com/haakonleg/imt2681-assig2/util"
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
		Ts:          util.NowMilli(),
		HDate:       igc.Date.String(),
		Pilot:       igc.Pilot,
		Glider:      igc.GliderType,
		GliderID:    igc.GliderID,
		TrackLength: util.CalTrackLen(igc.Points),
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
