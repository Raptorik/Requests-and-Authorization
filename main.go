package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"handlers/core"
	"log"
	"net/http"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}
func main() {

	// Simply returns a link to the login route
	http.HandleFunc("/", core.RootHandler)

	// Login route
	http.HandleFunc("/login/github/", core.GitGHubLoginHandler)

	// Github callback
	http.HandleFunc("/login/github/callback", core.GitHubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		core.LoggedInHandler(w, r, "")
	})

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}
