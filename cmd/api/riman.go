package main

import (
	"github.com/go-rod/rod"
	"time"
)

func login(loginUrl string, username string, password string) {
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(loginUrl)

	page.MustElement("div.static-menu-item").MustClick()

	page.MustElement("#mat-input-0").MustInput(username)

	page.MustElement("#mat-input-1").MustInput(password)

	page.MustElement(`[type="submit"]`).MustClick()

	/*	page.MustWaitStable().MustScreenshot("a.png") */

	time.Sleep(time.Hour)
}
