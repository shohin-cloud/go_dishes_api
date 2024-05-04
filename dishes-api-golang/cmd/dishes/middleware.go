package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetMember(r, model.AnonymousMember)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]
		v := validator.New()
		if model.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		member, err := app.models.Members.GetForToken(model.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		r = app.contextSetMember(r, member)
		next.ServeHTTP(w, r)
	})
}
