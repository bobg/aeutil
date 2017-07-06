package aeutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine/log"
)

func ReturnJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	err := e.Encode(v)
	if err != nil {
		HTTPErr(ctx, w, http.StatusInternalServerError, "writing response: %s", err)
		return
	}
}

func HTTPErr(ctx context.Context, w http.ResponseWriter, code int, f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	log.Errorf(ctx, msg)
	http.Error(w, msg, code)
}
