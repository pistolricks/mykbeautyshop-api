package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
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

	login(loginUrl, username, password)
}
