package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func (app *application) RimanLogin(loginUrl string, rimanStoreName string, username string, password string) (*rod.Page, *rod.Browser, []*proto.NetworkCookie) {
	// --allow-third-party-cookies

	if app.browser != nil {
		p := app.browser
		p.MustClose()
	}

	path, _ := launcher.LookPath()

	// homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	u := launcher.New().
		Headless(false).
		Devtools(true).
		NoSandbox(true).
		Bin(path).
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect().DefaultDevice(devices.LaptopWithHiDPIScreen)

	page := browser.MustPage(loginUrl)

	page.MustElement("div.static-menu-item").MustClick()
	page.MustElement("#mat-input-0").MustInput(username)
	page.MustElement("#mat-input-1").MustInput(password)
	page.MustElement(`[type="submit"]`).MustClick()

	cookies := browser.MustGetCookies()
	return page, browser, cookies
}

func (app *application) HomePage(rimanStoreName string, page *rod.Page, browser *rod.Browser, cookies []*proto.NetworkCookie) (*rod.Page, *rod.Browser, []*proto.NetworkCookie, error) {
	// networkCookie := networkCookies(cookies)

	homeUrl := fmt.Sprintf("https://mall.riman.com/%s/home", rimanStoreName)

	// page.MustSetCookies(networkCookie...)

	wait := page.MustWaitNavigation()
	page.MustNavigate(homeUrl)
	wait()

	app.page = page
	app.browser = browser

	newCookies, err := browser.GetCookies()
	if err != nil {
		fmt.Println(err)
		app.cookies = newCookies
		return page, browser, newCookies, err
	}

	app.cookies = newCookies

	return page, browser, newCookies, err
}

func (app *application) ProcessOrders(rimanStoreName string, page *rod.Page, browser *rod.Browser, cookies []*proto.NetworkCookie, orders []goshopify.Order) {
	orderCount := len(orders)

	switch orderCount := orderCount; {
	case orderCount == 1:
		app.SubmitOrder(rimanStoreName, page, browser, cookies, orders[0])
	case orderCount > 1:
		for _, order := range orders {
			app.SubmitOrder(rimanStoreName, page, browser, cookies, order)
		}
	}
}

func (app *application) SubmitOrder(rimanStoreName string, page *rod.Page, browser *rod.Browser, cookies []*proto.NetworkCookie, order goshopify.Order) {

	count := len(order.LineItems)

	for i, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/%s/products/%s", rimanStoreName, product.SKU)

		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl) // := browser.MustPage(productUrl)
		wait()

		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))
		page.MustElement("button.add-to-bag-btn").MustClick()
		page.MustWaitStable()

		fmt.Println("QTY, and Index")
		println(i + 1)
		println(count)
		/*
			switch {
			case i < count-1:
				page.MustElement("div.cart-btn").MustClick()
			case i == count-1:

				cookies, err := browser.GetCookies()
				if err != nil {
					return
				}

				app.processShipping(browser, page, cookies, order)

			}
		*/

		app.processShipping(browser, page, cookies, order)
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

func (app *application) processShipping(browser *rod.Browser, page *rod.Page, cookies []*proto.NetworkCookie, order goshopify.Order) {

	app.background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				app.insertShippingInfo(browser, page, checkoutUrl, order)

			default:
				fmt.Println("not right cookie")
			}

		}
	})

}

type StateObject = struct {
	Code  string
	Name  string
	name2 any
}

/* TODO: REMOVE HARD CODED EMAIL */

func (app *application) insertShippingInfo(browser *rod.Browser, page *rod.Page, checkoutUrl string, order goshopify.Order) {

	/*
		p := browser.MustPage(checkoutUrl)

		newCookies, err := browser.GetCookies()
		if err != nil {
			fmt.Println(err)
			return
		}

		networkCookie := networkCookies(newCookies)

		p.MustSetCookies(networkCookie...)

		wait := p.MustWaitNavigation()
		p.MustNavigate(checkoutUrl)
		wait()
	*/

	page.MustNavigate(checkoutUrl)

	shippingAddress := order.ShippingAddress

	firstName := strings.TrimSpace(shippingAddress.FirstName)
	lastName := strings.TrimSpace(shippingAddress.LastName)

	address1 := strings.TrimSpace(shippingAddress.Address1)
	address2 := strings.TrimSpace(shippingAddress.Address2)
	company := strings.TrimSpace(shippingAddress.Company)
	city := strings.TrimSpace(shippingAddress.City)
	province := strings.TrimSpace(shippingAddress.Province)
	provinceCode := strings.TrimSpace(shippingAddress.ProvinceCode)
	shortZip := strings.TrimSpace(shippingAddress.Zip[:5])
	zip := strings.TrimSpace(shippingAddress.Zip)

	phone := strings.Replace(strings.TrimSpace(shippingAddress.Phone), "+1", "", 1)
	// email := strings.TrimSpace(order.Email)

	page.MustElement("#firstName0").MustSelectAllText().MustInput(firstName)
	page.MustElement("#lastName0").MustSelectAllText().MustInput(lastName)

	removedAddress2 := strings.Replace(address1, address2, "", 1)
	removedCity := strings.Replace(removedAddress2, city, "", 1)
	removedProvince := strings.Replace(removedCity, province, "", 1)
	removedProvinceCode := strings.Replace(removedProvince, provinceCode, "", 1)
	removedZip := strings.Replace(removedProvinceCode, zip, "", 1)
	lineAddress := strings.Replace(removedZip, shortZip, "", 1)

	formattedAddress := strings.TrimSpace(lineAddress)

	address := fmt.Sprintf("%s %s, %s", formattedAddress, address2, shortZip)

	page.MustElement("#address10").MustSelectAllText().MustInput(address)
	page.MustElement("#address20").MustSelectAllText().MustInput(company)

	page.MustElement("#city0").MustSelectAllText().MustInput(city)
	// page.MustElement("#state0").MustSelect(provinceCode)
	page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip)

	page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)
	email := os.Getenv("ACCOUNT_EMAIL")
	page.MustElement("#email0").MustSelectAllText().MustInput(email)

	/* Need to add Province/State */
	// page.MustElement("#state0").MustSelectAllText().MustInput(province)
}

func (app *application) RimanLogout() bool {
	// --allow-third-party-cookies

	if app.browser != nil {
		p := app.browser
		p.MustClose()
	}

	return true
}
