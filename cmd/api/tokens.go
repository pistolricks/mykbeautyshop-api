package main

import (
	"errors"
	"fmt"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/validator"
	"net/http"
	"os"
	"time"
)

func (app *application) createRimanTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		RimanStoreName string `json:"rimanStoreName"`
		UserName       string `json:"userName"`
		Password       string `json:"password"`
		LoginUrl       string `json:"loginUrl"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	credentials := riman.Credentials{
		UserName: input.UserName,
		Password: input.Password,
	}

	v := validator.New()
	data.ValidatePasswordPlaintext(v, credentials.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.envars.Username = credentials.UserName
	app.envars.Password = credentials.Password
	app.envars.LoginUrl = input.LoginUrl

	err = os.Setenv("RIMAN_STORE_NAME", input.RimanStoreName)
	err = os.Setenv("LOGIN_URL", input.LoginUrl)
	err = os.Setenv("USERNAME", credentials.UserName)
	err = os.Setenv("Password", credentials.Password)

	post, err := riman.Login(credentials)

	err = app.writeJSON(w, http.StatusOK, envelope{"auth": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserName string `json:"userName"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	fmt.Println("USER")
	fmt.Println(user)

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	credentials := riman.Credentials{
		UserName: user.UserName,
		Password: input.Password,
	}

	fmt.Println("CREDENTIALS")
	fmt.Println(credentials)

	res, err := riman.Login(credentials)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println("RIMAN LOGIN")
	fmt.Println(res)

	token, err := app.models.Tokens.NewRid(user.ID, 24*time.Hour, data.ScopeAuthentication, res.Jwt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println("TOKEN")
	fmt.Println(token)

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user, "authentication_token": token, "auth": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		v.AddError("email", "user account must be activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}

		err = app.mailer.Send(user.Email, "token_password_reset.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	env := envelope{"message": "an email will be sent to you containing password reset instructions"}

	err = app.writeJSON(w, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Activated {
		v.AddError("email", "user has already been activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
		}

		err = app.mailer.Send(user.Email, "token_activation.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	env := envelope{"message": "an email will be sent to you containing activation instructions"}

	err = app.writeJSON(w, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
