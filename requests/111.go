package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	requests := []string{"Get", "Post"}

	for _, reques := range requests {
		switch reques {
		case "Get":
			var err error
			res, err := http.Get("https://api.github.com/users/Raptorik")
			if err != nil {
				log.Fatal(err)
			}
			log.Println("StatusCode:", res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Body: %s\n", body)
			type jsonUser struct {
				Name string `json:"login"`
			}
			user := jsonUser{}
			err = json.Unmarshal(body, &user)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Received user %s", user.Name)
		case "Post":
			c := http.Client{Timeout: time.Second}
			myJson := bytes.NewBuffer([]byte(`{"name": "Roman"}`))
			resp, err := c.Post(`https://api.github.com/users/Raptorik`, `application/json`, myJson)
			if err != nil {
				fmt.Printf("error: %s\\n", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			fmt.Printf("Body: %s\n", body)
		}

	}
}
