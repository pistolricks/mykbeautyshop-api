package riman

import (
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/go-rod/rod"
	"strconv"
)

func ProcessProducts(page *rod.Page, products []goshopify.LineItem) {

	for i, v := range products {
		productUrl := fmt.Sprintf("https://mall.riman.com/Werekbeauty/products/%s", products[i].SKU)
		wait := page.MustWaitNavigation()
		page.MustNavigate(productUrl)
		wait()
		fmt.Println(v.SKU)
		page.MustElement("input.quantity-input").MustSelectAllText().MustInput(strconv.Itoa(products[i].Quantity))
		page.MustWaitStable()
		page.MustElement("button.add-to-bag-btn").MustClick()
	}
}
