package riman

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Post struct {
	SecurityRedirect bool   `json:"securityRedirect"`
	Status           string `json:"-"`
	LiToken          string `json:"liToken"`
	LiUser           string `json:"liUser"`
	Jwt              string `json:"jwt"`
}

const apiUrl = "https://security-api.riman.com/api/v2/CheckAttemptsAndLogin"

func Login(credentials Credentials) (Post, error) {

	params := url.Values{}
	params.Add("userName", credentials.UserName)
	params.Add("password", credentials.Password)

	resp, err := http.PostForm(apiUrl, params)

	if err != nil {
		log.Printf("Posting failed: %s", err.Error())
		return Post{}, nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)
	// Unmarshal result

	post := Post{}

	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err.Error())
		return Post{}, nil
	}

	return post, err

}
