package main

import (
	"context"
	"fmt"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/pistolricks/kbeauty-api/internal/data"
	"net/http"
)

/* ORDER STATUS */
// - open
// - closed
// - cancelled
// - any

/* ORDER FULFILLMENT STATUS */
// - shipped
// - partial
// - unshipped
// - any
// - unfulfilled
// - fulfilled

/* ORDER FINANCIAL STATUS */
// authorized
// pending
// paid
// partially_paid
// refunded
// voided
// partially_refunded
// any
// unpaid

func (app *application) listOrdersByStatusHandler(w http.ResponseWriter, r *http.Request) {
	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := app.readStringParam("status", r)
	options := struct {
		Status string `url:"status"`
	}{s}

	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listOrdersByAllStatusValuesHandler(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	s := app.readStringParam("status", r)
	f := app.readStringParam("fulfillment_status", r)

	options := struct {
		Status            string `url:"status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{s, f}

	count, err := client.Order.Count(context.Background(), options)

	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) processOrders(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	options := struct {
		Status            string `url:"status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{"open", "unfulfilled"}

	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.background(func() {
		app.ProcessOrders(app.envars.LoginUrl, app.envars.Username, app.envars.Password, orders)
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listOrdersHandler(w http.ResponseWriter, r *http.Request) {

	shopApp := goshopify.App{
		ApiKey:      app.envars.ShopifyKey,
		ApiSecret:   app.envars.ShopifySecret,
		RedirectUrl: "https://example.com/callback",
		Scope:       "read_orders",
	}

	client, err := goshopify.NewClient(shopApp, app.envars.StoreName, app.envars.ShopifyToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	var input struct {
		Status            string
		FinancialStatus   string
		FulfillmentStatus string
		data.Filters
	}

	qs := r.URL.Query()

	fmt.Println(qs)

	input.Status = app.readString(qs, "status", "any")
	input.FinancialStatus = app.readString(qs, "financial_status", "any")
	input.FulfillmentStatus = app.readString(qs, "fulfillment_status", "any")

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id"}

	options := struct {
		Status            string `url:"status"`
		FinancialStatus   string `url:"financial_status"`
		FulfillmentStatus string `url:"fulfillment_status"`
	}{input.Status, input.FinancialStatus, input.FulfillmentStatus}

	fmt.Println(options)
	count, err := client.Order.Count(context.Background(), options)
	if err != nil {
		fmt.Println(err)
		return
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "count": count}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
