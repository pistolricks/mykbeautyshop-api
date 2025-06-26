package main

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

func setup() {
	app := goshopify.App{
		ApiKey:      "abcd",
		ApiSecret:   "efgh",
		RedirectUrl: "https://example.com/shopify/callback",
		Scope:       "read_products",
	}

	client, err := goshopify.NewClient(app, "shopname", "token")

	println(client, err)

}
