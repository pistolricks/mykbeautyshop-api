package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/products", app.listProductsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/orders", app.listOrdersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/orders/:status", app.listOrdersByStatusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/orders/:status/status/:fulfillment_status", app.listOrdersByAllStatusValuesHandler)

	router.HandlerFunc(http.MethodGet, "/v1/process/orders", app.processOrders)
	router.HandlerFunc(http.MethodPost, "/v1/process/order", app.processOrder)

	router.HandlerFunc(http.MethodGet, "/v1/riman/login", app.RimanLoginHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/orders", app.shippingHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/tracking", app.trackingHandler)

	/*
		router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

		router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	*/

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
