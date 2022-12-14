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

const clientID = "<7fed1fa2614ad27b52e5>"
const clientSecret = "<401052219a5e1a0f51da2b829a5419fc42167164>"

func LoggedInHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	// Set return type JSON
	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github/">Raptorik</a>`)
}

func GetGithubClientID() string {
	GitHubClientID, exists := os.LookupEnv("7fed1fa2614ad27b52e5")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	return GitHubClientID
}

func GetGithubClientSecret() string {
	GetGithubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	return GetGithubClientSecret
}
func GitGHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the environment variable
	GitHubClientID := GetGitHubClientID()

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		GetGitHubAccessToken,
		"http://localhost:3000/login/github/callback",
	)

	http.Redirect(w, r, redirectURL, 301)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	GitHubAccessToken := GetGitHubAccessToken(code)
	GitHubData := GetGitHubData(GitHubAccessToken)
	LoggedInHandler(w, r, GitHubData)
}
func GetGitHubAccessToken(code string) string {
	clientID := GetGitHubClientID()
	clientSecret := GetGitHubClientSecret()
	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"7fed1fa2614ad27b52e5":                     clientID,
		"401052219a5e1a0f51da2b829a5419fc42167164": clientSecret,
		"code": code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Response body converted to string if ied JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from GitHub
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert string if ied JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func GetGitHubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
