package util

import "math"

const earthRadius = 6371000 // 地球半徑，單位為米

type LatLng struct {
	Lat float64
	Lng float64
}

func (l LatLng) IsWithinRadiusOf(latLng LatLng, radius float64) bool {
	return l.DistanceTo(latLng) <= radius
}

func (l LatLng) DistanceTo(newL LatLng) float64 {
	var degToRad = func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	lat1Rad := degToRad(l.Lat)
	lng1Rad := degToRad(l.Lng)
	lat2Rad := degToRad(newL.Lat)
	lng2Rad := degToRad(newL.Lng)

	// 計算兩點之間的距離
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
