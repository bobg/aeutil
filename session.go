package aeutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Session struct {
	Ident string

	UserKey *User

	Closed bool

	CreatedWhen  time.Time
	CreatedWhere appengine.GeoPoint
}

type SessionStore struct {
	kind   string
	cookie string
}

func NewSessionStore(kind, cookie string) *SessionStore {
	return &SessionStore{kind: kind, cookie: cookie}
}

func (ss *SessionStore) Get(ctx context.Context, r *http.Request) (*datastore.Key, *Session, bool, error) {
	newSession := func() (*datastore.Key, *Session, bool, error) {
		k, s, err := ss.New(ctx)
		return k, s, true, err
	}
	cookie, err := r.Cookie(ss.cookie)
	if err == http.ErrNoCookie {
		return newSession()
	}
	if err != nil {
		return nil, nil, false, err
	}
	k := datastore.NewKey(ctx, ss.kind, cookie.Value, 0, nil)
	var s Session
	err = datastore.Get(ctx, k, &s)
	if err != nil || s.Closed {
		return newSession()
	}
	return k, &s, false, nil
}

func (ss *SessionStore) New(ctx context.Context) (*datastore.Key, *Session, error) {
	var ident [32]byte
	_, err := rand.Read(ident[:])
	if err != nil {
		return nil, nil, err
	}
	identHex := hex.EncodeToString(ident)
	geoPoint, _ := GetGeoPoint(r)
	s := &Session{
		Ident:        identHex,
		CreatedWhen:  time.Now(),
		CreatedWhere: geoPoint,
	}
	k := datastore.NewKey(ctx, ss.kind, identHex, 0, nil)
	_, err := datastore.Put(ctx, k, s)
	return k, s, err
}

func (ss *SessionStore) WriteCookie(w http.ResponseWriter, k *datastore.Key) {
	cookie := &http.Cookie{
		Name:  ss.cookie,
		Value: k.StringID(),
	}
	http.SetCookie(w, cookie)
}
