package main

import (
	"github.com/go-rod/rod"
)

func (app *application) login() *rod.Page {
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(app.envars.LoginUrl)

	page.MustElement("div.static-menu-item").MustClick()

	page.MustElement("#mat-input-0").MustInput(app.envars.Username)

	page.MustElement("#mat-input-1").MustInput(app.envars.Password)

	page.MustElement(`[type="submit"]`).MustClick()

	/*	page.MustWaitStable().MustScreenshot("a.png") */

	// time.Sleep(time.Hour)

	return page
}
