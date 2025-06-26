package main

import (
	"context"
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func (app *application) connect() {

	// redirectUrl := fmt.Sprintf("localhost:4000/%s/callback", app.envars.StoreName)

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_products, read_product_listings",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	numProducts, err := client.Product.Count(context.Background(), nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// productListings, err := client.Product.List(context.Background(), nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(numProducts)

	// fmt.Println(productListings)

}
