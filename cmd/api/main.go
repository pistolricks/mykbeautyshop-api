package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

type Envars struct {
	StoreName     string
	LoginUrl      string
	Username      string
	Password      string
	ShopifyToken  string
	ShopifyKey    string
	ShopifySecret string
}

type application struct {
	logger *slog.Logger
	envars *Envars
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	storeName := os.Getenv("STORE_NAME")
	if storeName == "" {
		fmt.Println("missing store name")
		return
	}

	loginUrl := os.Getenv("LOGIN_URL")
	if loginUrl == "" {
		fmt.Println("missing login url")
		return
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		fmt.Println("missing username")
		return
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		fmt.Println("missing password")
		return
	}

	shopifyToken := os.Getenv("SHOPIFY_TOKEN")
	if shopifyToken == "" {
		fmt.Println("missing shop token")
		return
	}

	shopifyKey := os.Getenv("SHOPIFY_KEY")
	if shopifyKey == "" {
		fmt.Println("missing shop key")
		return
	}

	shopifySecret := os.Getenv("SHOPIFY_SECRET")
	if shopifySecret == "" {
		fmt.Println("missing shop secret")
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	vars := &Envars{StoreName: storeName, LoginUrl: loginUrl, Username: username, Password: password, ShopifyToken: shopifyToken, ShopifyKey: shopifyKey, ShopifySecret: shopifySecret}

	fmt.Println(vars)

	app := &application{
		logger: logger,
		envars: vars,
	}

	app.setup()

	// app.login()
}
