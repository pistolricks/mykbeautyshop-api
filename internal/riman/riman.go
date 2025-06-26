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

		AddToCart(page, strconv.Itoa(products[i].Quantity))

		fmt.Printf("2+%d =", i)

	}
}

func OpenProductPage(page *rod.Page, productID string) {

	productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", productID)

	page.MustNavigate(productUrl).MustWaitStable()

}

func AddToCart(page *rod.Page, quantity string) {

	page.MustElement("input.quantity-input").MustInput(quantity)

	page.MustElement("button.add-to-bag-btn").MustClick()

	/*	page.MustWaitStable().MustScreenshot("a.png") */

}
