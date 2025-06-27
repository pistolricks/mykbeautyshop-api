package main

import (
	"fmt"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Login(loginUrl string, username string, password string) *rod.Page {
	browser := rod.New().MustConnect().NoDefaultDevice()

	page := browser.MustPage(loginUrl)

	page.MustElement("div.static-menu-item").MustClick()
	page.MustElement("#mat-input-0").MustInput(username)
	page.MustElement("#mat-input-1").MustInput(password)
	page.MustElement(`[type="submit"]`).MustClick()

	/*	page.MustWaitStable().MustScreenshot("a.png") */
	// time.Sleep(time.Hour)

	return page
}

func (app *application) ProcessOrders(loginUrl string, username string, password string, orders []goshopify.Order) {
	orderCount := len(orders)

	switch orderCount := orderCount; {
	case orderCount == 1:
		app.SubmitOrder(loginUrl, username, password, orders[0])
	case orderCount > 1:
		for _, order := range orders {
			app.SubmitOrder(loginUrl, username, password, order)
		}
	}
}

func (app *application) SubmitOrder(loginUrl string, username string, password string, order goshopify.Order) {

	browser := rod.New().MustConnect().NoDefaultDevice()

	page := browser.MustPage(loginUrl)

	page.MustElement("div.static-menu-item").MustClick()
	page.MustElement("#mat-input-0").MustInput(username)
	page.MustElement("#mat-input-1").MustInput(password)
	page.MustElement(`[type="submit"]`).MustClick()

	cookies := browser.MustGetCookies()

	networkCookie := networkCookies(cookies)

	count := len(order.LineItems)

	for i, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", product.SKU)

		page := browser.MustPage(productUrl)

		page.MustSetCookies(networkCookie...)

		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl)
		wait()

		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))
		page.MustElement("button.add-to-bag-btn").MustClick()
		page.MustWaitStable()

		switch {
		case i < count-1:
			page.MustElement("div.cart-btn").MustClick()
		case i == count-1:

			cookies, err := browser.GetCookies()
			if err != nil {
				return
			}

			app.cookieServer(page, cookies)

		}

	}
}

func networkCookies(cookies []*proto.NetworkCookie) []*proto.NetworkCookieParam {

	var networkCookie []*proto.NetworkCookieParam

	for _, cookie := range cookies {
		networkCookie = append(networkCookie, &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
			SameSite: cookie.SameSite,
			Expires:  cookie.Expires,
		})
	}

	return networkCookie
}

func (app *application) cookieServer(page *rod.Page, cookies []*proto.NetworkCookie) {

	app.background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				shippingInfo(page, checkoutUrl)

			default:
				fmt.Println("not right cookie")
			}

		}
	})

}

func shippingInfo(page *rod.Page, checkoutUrl string) {
	wait := page.MustWaitNavigation()
	page.MustNavigate(checkoutUrl)
	wait()
}
