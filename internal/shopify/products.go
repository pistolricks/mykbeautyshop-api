package shopify

import (
	"context"
	"fmt"
	"github.com/vinhluan/go-shopify-graphql/model"
)

func ListAllProducts() ([]*model.Collection, error) {

	// Get all collections
	collections, err := Client().Collection.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	// Print out the result
	for _, c := range collections {
		fmt.Println(c.Handle)
	}

	return collections, err
}
