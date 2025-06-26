package riman

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"strconv"
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

func SubmitOrder(loginUrl string, username string, password string, order goshopify.Order) {

	page := Login(loginUrl, username, password)

	for index, product := range order.LineItems {
		productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", product.SKU)
		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl)
		wait()
		fmt.Println(order.LineItems[index])
		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(product.Quantity))
		page.MustWaitStable()
		// p.MustElement("button.add-to-bag-btn").MustElement(`[type="button"]`).MustClick()
	}
}
