package data

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func OrderUpdateTracking(orderId string, token string) string {

	path := fmt.Sprintf("/api/v1/orders/%s/shipment-products", orderId)

	u := &url.URL{
		Scheme: "https",
		Host:   "cart-api.riman.com",
		Path:   path,
	}

	q := u.Query()
	q.Add("token", token)

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(u.String())
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	bodyString := string(resBody)

	return bodyString

}

func OrderUpdateFulfillment() {

}
