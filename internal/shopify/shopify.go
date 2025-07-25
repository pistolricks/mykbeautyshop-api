package shopify

import (
	shopApi "github.com/vinhluan/go-shopify-graphql"
	"os"
)

func Client() *shopApi.Client {
	c := shopApi.NewDefaultClient()

	// Or if you are a fan of options
	c = shopApi.NewClient(os.Getenv("STORE_NAME"),
		shopApi.WithToken(os.Getenv("STORE_PASSWORD")),
		shopApi.WithVersion("2023-07"),
		shopApi.WithRetries(5))

	return c
}
