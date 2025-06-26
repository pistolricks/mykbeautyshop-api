package riman

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"strconv"
)

func ProcessProducts(page *rod.Page, products []goshopify.LineItem) {
	for i, v := range products {

		wait := page.MustWaitNavigation()

		OpenProductPage(page, v.SKU)

		wait()

		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(products[i].Quantity))
		page.MustWaitStable()
		AddToCart(page)

		fmt.Printf("2+%d =", i)

	}
}

func OpenProductPage(page *rod.Page, productID string) {

	productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", productID)

	page.MustNavigate(productUrl).MustWaitStable()

}

func AddToCart(page *rod.Page) {

	page.MustElement("button.add-to-bag-btn").MustClick()

	/*	page.MustWaitStable().MustScreenshot("a.png") */

}
