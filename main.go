package main

import (
	"fmt"
	"github.com/Raptorik/oAuth/tree/main/mygithubAUTH/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}
func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/login/github/", handlers.GitHubLoginHandler)
	http.HandleFunc("/login/github/callback", handlers.GitHubCallbackHandler)
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoggedInHandler(w, r, "")
	})
	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}
