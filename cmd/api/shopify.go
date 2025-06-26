package main

import (
	"context"
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func (app *application) setup() {

	// redirectUrl := fmt.Sprintf("localhost:4000/%s/callback", app.envars.StoreName)

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_products",
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

	fmt.Println(numProducts)

}
