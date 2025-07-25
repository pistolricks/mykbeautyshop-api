package main

import (
	"expvar"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/pistolricks/kbeauty-api/graph"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
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

	router.HandlerFunc(http.MethodPost, "/v1/attributes/metafield", app.updateOrderMetaField)

	router.HandlerFunc(http.MethodGet, "/v1/process/orders", app.processOrders)
	router.HandlerFunc(http.MethodPost, "/v1/process/order", app.processOrder)

	router.HandlerFunc(http.MethodGet, "/v1/riman/orders", app.shippingHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/tracking", app.trackingHandler)

	router.HandlerFunc(http.MethodPost, "/v1/riman/login", app.createRimanTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.createAuthenticationTokenHandler)
	/*

		router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)


		router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	*/

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/v2/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/v2/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "4000")

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
