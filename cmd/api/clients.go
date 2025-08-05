package main

import (
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"github.com/pistolricks/kbeauty-api/internal/riman"
	"github.com/pistolricks/kbeauty-api/internal/validator"
	"net/http"
	"os"
	"slices"
)

func (app *application) homePageHandler(w http.ResponseWriter, r *http.Request) {

	rimanStoreName := os.Getenv("RIMAN_STORE_NAME")
	rimanRid := os.Getenv("USERNAME")
	rimanUrl := os.Getenv("LOGIN_URL")
	currentBrowser := app.browser
	currentCookies := app.cookies

	idx := slices.IndexFunc(currentCookies, func(c *proto.NetworkCookie) bool { return c.Name == "token" })

	var foundCookie proto.NetworkCookie

	if idx != -1 {
		foundCookie := currentCookies[idx]
		fmt.Printf("Found client: %v\n", foundCookie)
		// Output: Found client: {ID:2 Name:Jane Smith}
	} else {
		fmt.Println("Client not found")
	}

	app.envars.RimanStoreName = rimanStoreName
	app.envars.LoginUrl = rimanUrl
	app.envars.Username = rimanRid
	app.envars.Token = foundCookie.Value

	page, browser, cookies, _ := app.HomePage(rimanStoreName, currentBrowser, currentCookies)

	fmt.Println(cookies)

	err := app.writeJSON(w, http.StatusOK, envelope{"page": page, "browser": browser, "cookies": cookies, "rid": rimanRid, "store": rimanStoreName, "url": rimanUrl}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

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

	err = os.Setenv("RIMAN_STORE_NAME", input.RimanStoreName)
	err = os.Setenv("LOGIN_URL", input.LoginUrl)
	err = os.Setenv("USERNAME", credentials.UserName)
	err = os.Setenv("Password", credentials.Password)

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
