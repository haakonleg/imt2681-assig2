package paragliding

import (
	"fmt"
	"net/http"
	"time"
)

var appStartTime = time.Now()

type apiInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// Send API info
func sendAPIInfo(req *Request) {
	info := &apiInfo{uptime(), "Service for Paragliding tracks.", "v1"}
	req.SendJSON(info, http.StatusOK)
}

// uptime returns the app uptime in ISO 8601 duration format
func uptime() string {
	// Seconds duration since app start
	duration := time.Since(appStartTime)

	sec := int(duration.Seconds()) % 60
	min := int(duration.Minutes()) % 60
	hour := int(duration.Hours()) % 24
	day := int(duration.Hours()/24) % 7
	month := int(duration.Hours()/24/7/4) % 12
	year := int(duration.Hours() / 24 / 365)

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", year, month, day, hour, min, sec)
}
