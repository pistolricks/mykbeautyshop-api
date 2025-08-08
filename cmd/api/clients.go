package main

import (
	"fmt"
	"net/http"

	"github.com/pistolricks/mykbeautyshop-api/internal/data"
	"github.com/pistolricks/mykbeautyshop-api/internal/riman"
	"github.com/pistolricks/mykbeautyshop-api/internal/validator"
)

func (app *application) findCookieValue() *string {
	for i := range app.cookies {
		if app.cookies[i].Name == "token" {
			app.envars.Token = app.cookies[i].Value
			fmt.Println("app.envars.Token")
			fmt.Println(app.envars.Token)
			return &app.cookies[i].Value
		}
	}
	// Return nil if no product is found
	return nil
}

func (app *application) clientLoginHandler(w http.ResponseWriter, r *http.Request) {

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

	app.envars.RimanStoreName = input.RimanStoreName
	app.envars.Username = credentials.UserName
	app.envars.Password = credentials.Password
	app.envars.LoginUrl = input.LoginUrl

	// err = os.Setenv("RIMAN_STORE_NAME", input.RimanStoreName)
	// err = os.Setenv("LOGIN_URL", input.LoginUrl)
	// err = os.Setenv("USERNAME", credentials.UserName)
	// err = os.Setenv("Password", credentials.Password)

	page, browser, cookies := app.RimanLogin(input.LoginUrl, input.RimanStoreName, credentials.UserName, credentials.Password)

	app.page = page
	app.browser = browser
	app.cookies = cookies
	fmt.Println(browser)

	err = app.writeJSON(w, http.StatusOK, envelope{"page": app.page, "browser": app.browser, "cookies": app.cookies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {

	rimanStoreName := app.envars.RimanStoreName // os.Getenv("RIMAN_STORE_NAME")
	rimanRid := app.envars.Username             // os.Getenv("USERNAME")
	rimanUrl := app.envars.LoginUrl             // os.Getenv("LOGIN_URL")
	currentPage := app.page
	currentBrowser := app.browser
	currentCookies := app.cookies

	page, browser, cookies, _ := app.HomePage(rimanStoreName, currentPage, currentBrowser, currentCookies)

	fmt.Println(cookies)

	token := app.findCookieValue()
	fmt.Println("TOKEN")
	fmt.Println(token)

	err := app.writeJSON(w, http.StatusOK, envelope{"page": page, "browser": browser, "cookies": cookies, "rid": rimanRid, "store": rimanStoreName, "url": rimanUrl}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listClientsHandler(w http.ResponseWriter, r *http.Request) {

	clients, metadata, err := app.models.Clients.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"clients": clients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) clientLogoutHandler(w http.ResponseWriter, r *http.Request) {

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

	if input.RimanStoreName == app.envars.RimanStoreName {
		app.envars.RimanStoreName = ""
	}
	if input.UserName == app.envars.Username {
		app.envars.Username = ""
	}
	if input.Password == app.envars.Password {
		app.envars.Password = ""
	}

	if input.LoginUrl == app.envars.LoginUrl {
		app.envars.LoginUrl = ""
	}

	isLoggedOut := app.RimanLogout()

	err = app.writeJSON(w, http.StatusOK, envelope{"logged_out": isLoggedOut}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
