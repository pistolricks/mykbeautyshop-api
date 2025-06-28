package main

import (
	"fmt"
	"github.com/go-rod/rod/lib/devices"
	"strconv"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Login(loginUrl string, username string, password string) *rod.Page {
	browser := rod.New().MustConnect().DefaultDevice(devices.LaptopWithHiDPIScreen)

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

	browser := rod.New().MustConnect().DefaultDevice(devices.LaptopWithHiDPIScreen)

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

			app.processShipping(page, cookies, order)

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

func (app *application) processShipping(page *rod.Page, cookies []*proto.NetworkCookie, order goshopify.Order) {

	app.background(func() {

		for _, cookie := range cookies {

			switch n := cookie.Name; n {
			case "cartKey":
				fmt.Println(cookie.Value)
				fmt.Println("it worked")

				cartValue := cookie.Value
				fmt.Println(cartValue)

				checkoutUrl := fmt.Sprintf("https://mall.riman.com/checkout/shipping?cartKey=%s", cartValue)

				shippingInfo(page, checkoutUrl, order)

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

func shippingInfo(page *rod.Page, checkoutUrl string, order goshopify.Order) {
	wait := page.MustWaitNavigation()
	page.MustNavigate(checkoutUrl)
	wait()

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
	email := strings.TrimSpace(order.Email)

	page.MustElement("#firstName0").MustSelectAllText().MustInput(firstName)
	page.MustElement("#lastName0").MustSelectAllText().MustInput(lastName)

	removedAddress2 := strings.Replace(address1, address2, "", 1)
	removedCity := strings.Replace(removedAddress2, city, "", 1)
	removedProvince := strings.Replace(removedCity, province, "", 1)
	removedProvinceCode := strings.Replace(removedProvince, provinceCode, "", 1)
	removedZip := strings.Replace(removedProvinceCode, zip, "", 1)
	formatted1 := strings.Replace(removedZip, shortZip, "", 1)

	address := fmt.Sprintf("%s %s, %s", formatted1, address2, zip)

	page.MustElement("#address10").MustSelectAllText().MustInput(address)
	page.MustElement("#address20").MustSelectAllText().MustInput(company)

	page.MustElement("#city0").MustSelectAllText().MustInput(city)
	// page.MustElement("#state0").MustSelect(provinceCode)
	page.MustElement("#postalCode0").MustSelectAllText().MustInput(zip)

	page.MustElement("#phoneNumber0").MustSelectAllText().MustInput(phone)
	page.MustElement("#email0").MustSelectAllText().MustInput(email)

	/* Need to add Province/State */
	// page.MustElement("#state0").MustSelectAllText().MustInput(province)
}
