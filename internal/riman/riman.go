package riman

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"strconv"
)

func ProcessOrders(loginUrl string, username string, password string, orders []goshopify.Order) {
	orderCount := len(orders)

	switch orderCount := orderCount; {
	case orderCount == 1:
		SubmitOrder(loginUrl, username, password, orders[0])
	case orderCount > 1:
		for _, order := range orders {
			SubmitOrder(loginUrl, username, password, order)
		}
	}
}

func SubmitOrder(loginUrl string, username string, password string, order goshopify.Order) {

	browser := rod.New().MustConnect().NoDefaultDevice()

	page := browser.MustPage(loginUrl)

	page.MustElement("div.static-menu-item").MustClick()
	page.MustElement("#mat-input-0").MustInput(username)
	page.MustElement("#mat-input-1").MustInput(password)
	page.MustElement(`[type="submit"]`).MustClick()

	for _, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", product.SKU)
		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl)
		wait()

		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))

		page.MustElement("button.add-to-bag-btn").MustClick()
		page.MustWaitStable()
		page.MustElement("div.cart-btn").MustClick()
	}
}
