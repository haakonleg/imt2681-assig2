package util

import (
	"math"
	"strconv"
	"time"

	igc "github.com/marni/goigc"
)

const (
	earthRadius = float64(6371)
)

// Get current UNIX timestamp in miliseconds
func NowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Calculate the track length, the sum of all distances between each point
func CalTrackLen(points []igc.Point) string {
	d := 0.0
	for i := 0; i < len(points)-1; i++ {
		d += haversine(points[i].Lng.Radians(), points[i].Lat.Radians(), points[i+1].Lng.Radians(), points[i+1].Lat.Radians())
	}

	return strconv.FormatFloat(d, 'f', 2, 64) + "km"
}

// Formula to calculate distance between two coordinates on earth sphere, in km
// The anlges are already in radians
func haversine(rLon1, rLat1, rLon2, rLat2 float64) float64 {
	dLat := rLat2 - rLat1
	dLon := rLon2 - rLon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
