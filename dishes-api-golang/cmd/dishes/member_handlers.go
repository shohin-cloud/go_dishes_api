package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

func (app *application) registerMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	member := &model.Member{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = member.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateMember(v, member); !v.Valid() {
		fmt.Println("ERROR HERE?")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Members.Insert(member)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a member with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForMember(member.ID, "posts:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(member.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var res struct {
		Token  *string       `json:"token"`
		Member *model.Member `json:"member"`
	}

	res.Token = &token.Plaintext
	res.Member = member

	app.writeJSON(w, http.StatusCreated, envelope{"member": res}, nil)
}

func (app *application) activateMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if model.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	member, err := app.models.Members.GetForToken(model.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	member.Activated = true

	err = app.models.Members.Update(member)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForMember(model.ScopeActivation, member.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
