package main

import (
	"context"
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"net/http"
)

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {

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

	count, err := client.Product.Count(context.Background(), nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	products, err := client.Product.List(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
