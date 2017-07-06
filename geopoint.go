package aeutil

import (
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/appengine"
)

func GetGeoPoint(r *http.Request) (geoPoint appengine.GeoPoint, err error) {
	s := r.Header.Get("X-AppEngine-CityLatLong") // silently ignore errors
	strs := strings.Split(s, ",")
	if len(strs) != 2 {
		err = ErrInvalidLatLong
		return
	}
	geoPoint.Lat, err = strconv.ParseFloat(strs[0], 64)
	if err != nil {
		return
	}
	geoPoint.Lng, err = strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return
	}
	if !geoPoint.Valid() {
		err = ErrInvalidGeoPoint
	}
	return
}
