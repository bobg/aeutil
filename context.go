package aeutil

import (
	"context"

	"google.golang.org/appengine/datastore"
)

type (
	sessionCtxKeyType string
	sessionCtxPair    struct {
		k *datastore.Key
		s *Session
	}
)

const sessionCtxKey = sessionCtxKeyType("session")

func SessionCtx(ctx context.Context, k *datastore.Key, s *Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey, sessionCtxPair{k, s})
}

func CtxSession(ctx context.Context) (*datastore.Key, *Session) {
	v := ctx.Value(sessionCtxKey)
	if p, ok := v.(sessionCtxPair); ok {
		return p.k, p.s
	}
	return nil, nil
}
