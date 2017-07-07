package aeutil

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/appengine"
)

var ErrInvalidGeoPoint = errors.New("invalid geopoint")

func GetGeoPoint(r *http.Request) (geoPoint appengine.GeoPoint, err error) {
	s := r.Header.Get("X-AppEngine-CityLatLong") // silently ignore errors
	strs := strings.Split(s, ",")
	if len(strs) != 2 {
		err = ErrInvalidGeoPoint
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
