package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	GitHubClientID := GetGitHubClientID()

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", GitHubClientID, "http://localhost:3000/login/github/callback")

	http.Redirect(w, r, redirectURL, 301)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := GetGitHubAccessToken(code)
	githubData := GetGithubData(githubAccessToken)
	LoggedInHandler(w, r, githubData)
}

func GetGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}

func GetGitHubAccessToken(code string) string {
	clientID := GetGitHubClientID()
	clientSecret := GetGitHubClientSecret()
	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)
	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := ioutil.ReadAll(resp.Body)

	type GitHubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp GitHubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)
	return ghresp.AccessToken
}

func GetGitHubClientID() string {
	GitHubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	return GitHubClientID
}

func GetGitHubClientSecret() string {
	GitHubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	return GitHubClientSecret
}
