package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var clientID = ""
var clientSecret = ""
var redirectURL = "https://localhost:18888"
var state = "your state"

func main() {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "gist"},
		Endpoint:     github.Endpoint,
	}

	var token *oauth2.Token
	file, err := os.Open("access_token.json")
	if err == nil {
		defer file.Close()
	}

	if os.IsNotExist(err) {
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)

		code := make(chan string)
		var server *http.Server
		server = &http.Server{
			Addr: ":18888",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, "<html><script>window.open('about:blank','_self').close()</script></html>")
				w.(http.Flusher).Flush()
				code <- r.URL.Query().Get("code")
				server.Shutdown(context.Background())
			}),
		}
		go server.ListenAndServe()
		open.Start(url)

		token, err = conf.Exchange(oauth2.NoContext, <-code)
		if err != nil {
			panic(err)
		}
		file, err := os.Create("access_token.json")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(file).Encode(token)

		fmt.Println("obtained access token.")
	} else if err == nil {
		token = &oauth2.Token{}
		json.NewDecoder(file).Decode(token)

		fmt.Println("access token had already been obtained.")
	} else {
		panic(err)
	}
}
