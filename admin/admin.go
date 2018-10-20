package admin

import (
	"net/http"
	"strconv"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
)

type AdminHandler struct {
	db *mdb.Database
}

func NewAdminHandler(db *mdb.Database) *AdminHandler {
	return &AdminHandler{
		db: db}
}

// GetTrackCount is a handler for GET /admin/api/tracks_count
// It returns the total number of registered tracks in the database
func (ah *AdminHandler) GetTrackCount(req *router.Request) {
	tCnt, err := ah.db.Count(mdb.TRACKS)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"})
		return
	}
	req.SendText(strconv.FormatInt(tCnt, 10), http.StatusOK)
}

// DeleteAllTracks is a handler for DELETE /admin/api/tracks
// It deletes all the registered tracks from the database
func (ah *AdminHandler) DeleteAllTracks(req *router.Request) {
	_, err := ah.db.Delete(mdb.TRACKS, nil)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"})
		return
	}
	req.SendText("Everything deleted", http.StatusOK)
}
