package aeutil

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
)

type User struct {
	Email            string
	Verified         bool
	VerificationCode int
	CreatedWhen      time.Time
}

type UserStore struct {
	kind string
}

var ErrNotVerified = errors.New("not verified")

func NewUserStore(kind string) *UserStore {
	return &UserStore{kind: kind}
}

func (us *UserStore) Get(ctx context.Context, email string) (*datastore.Key, *User, bool, error) {
	email = canonicalizeEmail(email)
	k := datastore.NewKey(ctx, us.kind, email, 0, nil)
	var u User
	err := datastore.Get(ctx, k, &u)
	if err == datastore.ErrNoSuchEntity {
		u.Email = email
		u.VerificationCode = newVerificationCode()
		u.CreatedWhen = time.Now()
		_, err = datastore.Put(ctx, k, &u)
		return k, &u, true, err
	}
	return k, &u, false, err
}

func (us *UserStore) Verify(ctx context.Context, k *datastore.Key, code int) error {
	var u User
	err := datastore.Get(ctx, k, &u)
	if err != nil {
		return err
	}
	if u.Verified {
		return nil
	}
	if code != u.VerificationCode {
		u.VerificationCode = newVerificationCode()
		datastore.Put(ctx, k, &u) // ignore errors
		return ErrNotVerified
	}
	u.Verified = true
	_, err = datastore.Put(ctx, k, &u)
	return err
}

func newVerificationCode() int {
	return 100000 + rand.Intn(900000)
}

func canonicalizeEmail(s string) string {
	// TODO: better canonicalization
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}
