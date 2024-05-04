package main

import (
	"context"
	"net/http"

	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
)

type contextKey string

const memberContextKey = contextKey("member")

func (app *application) contextSetMember(r *http.Request, member *model.Member) *http.Request {
	ctx := context.WithValue(r.Context(), memberContextKey, member)
	return r.WithContext(ctx)
}

func (app *application) contextGetMember(r *http.Request) *model.Member {
	member, ok := r.Context().Value(memberContextKey).(*model.Member)
	if !ok {
		panic("missing member value in request context")
	}
	return member
}
