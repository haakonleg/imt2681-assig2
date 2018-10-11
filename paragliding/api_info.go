package paragliding

import (
	"fmt"
	"net/http"
	"time"

	"github.com/haakonleg/imt2681-assig2/router"
)

type ApiInfoHandler struct {
	startTime time.Time
}

func NewInfoHandler() *ApiInfoHandler {
	return &ApiInfoHandler{
		startTime: time.Now()}
}

// Send API info
func (aih *ApiInfoHandler) getAPIInfo(req *router.Request) {
	response := struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}{
		Uptime:  uptime(&aih.startTime),
		Info:    "Service for Paragliding tracks.",
		Version: "v1"}

	req.SendJSON(&response, http.StatusOK)
}

// uptime returns the app uptime in ISO 8601 duration format
func uptime(startTime *time.Time) string {
	// Seconds duration since app start
	duration := time.Since(*startTime)

	sec := int(duration.Seconds()) % 60
	min := int(duration.Minutes()) % 60
	hour := int(duration.Hours()) % 24
	day := int(duration.Hours()/24) % 7
	month := int(duration.Hours()/24/7/4) % 12
	year := int(duration.Hours() / 24 / 365)

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", year, month, day, hour, min, sec)
}
