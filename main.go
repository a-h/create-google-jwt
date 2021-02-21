package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/browser"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var flagBaseURL = flag.String("baseURL", "localhost:9999", "Set the URL that the local Web server will listen on.")
var flagClientID = flag.String("clientID", "", "The Client ID, e.g. xxxx-yyyy.apps.googleusercontent.com")
var flagClientSecret = flag.String("clientSecret", "", "The Client secret.")
var flagScopes = flag.String("scopes", "https://www.googleapis.com/auth/userinfo.email", "A comma separated list of scopes to request.")

func main() {
	flag.Parse()
	if *flagClientID == "" || *flagClientSecret == "" {
		fmt.Println("Please provide the clientID and clientSecret flags.")
		flag.Usage()
		os.Exit(1)
	}
	conf := &oauth2.Config{
		RedirectURL:  "http://" + *flagBaseURL + "/Callback",
		ClientID:     *flagClientID,
		ClientSecret: *flagClientSecret, // Get you own value for this from https://console.cloud.google.com/apis/api/cloudidentity.googleapis.com.
		Scopes:       strings.Split(*flagScopes, ","),
		Endpoint:     google.Endpoint,
	}
	state := uuid.NewV4().String()

	tokensChan := make(chan string, 1)
	s := http.Server{
		Addr: *flagBaseURL,
	}
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if actualState := r.URL.Query().Get("state"); actualState != state {
			http.Error(w, "unexpected authentication value", http.StatusUnauthorized)
			return
		}
		tokensChan <- r.URL.Query().Get("code")
		close(tokensChan)
		w.Write([]byte("Authentication complete, you can close this window."))
	})
	go func() {
		err := s.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("error starting web server: %v", err)
		}
	}()
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	ctx := context.Background()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Opening auth URL: %v\n", url)
	browser.OpenURL(url)
	code := <-tokensChan
	s.Close()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Printf("Access Token: %s %s\n", tok.Type(), tok.AccessToken)
	fmt.Println()
	idToken := tok.Extra("id_token").(string)
	fmt.Printf("Authorization: %s\n", idToken)
}
