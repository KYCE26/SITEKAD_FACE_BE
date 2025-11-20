package helper

import "math"

const earthRadiusMeters = 6371000 
func Geolocation(lat1, lon1, lat2, lon2 float64) float64 {
	latRad1 := lat1 * math.Pi / 180
	lonRad1 := lon1 * math.Pi / 180
	latRad2 := lat2 * math.Pi / 180
	lonRad2 := lon2 * math.Pi / 180

	diffLat := latRad2 - latRad1
	diffLon := lonRad2 - lonRad1

	a := math.Sin(diffLat/2)*math.Sin(diffLat/2) +
		math.Cos(latRad1)*math.Cos(latRad2)*
			math.Sin(diffLon/2)*math.Sin(diffLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}