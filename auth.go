package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// generate random string for state
func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// setup authentication
func SetUpAuth(config *oauth2.Config, provider *oidc.Provider) {
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		state, err := generateState()
		if err != nil {
			http.Error(w, "failed to generate state", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "oauthstate",
			Value: state,
		})
		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusSeeOther)
	})
	http.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		oauthstate, _ := r.Cookie("oauthstate")
		if r.FormValue("state") != oauthstate.Value {
			http.Error(w, "invalid oauth state", http.StatusBadRequest)
			return
		}
		token, err := config.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			http.Error(w, "failed to exchange token", http.StatusInternalServerError)
			return
		}
		//get the id token
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "no id_token", http.StatusInternalServerError)
			return
		}
		verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")})
		//verify the id token
		idToken, err := verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, "failed to verify ID Token", http.StatusInternalServerError)
			return
		}

		//get userInfo
		userInfo := map[string]string{}
		if err := idToken.Claims(&userInfo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(userInfo)

		//set the token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token.AccessToken,
		})

		//redirect to the graph page
		w.WriteHeader(http.StatusPermanentRedirect)
		http.Redirect(w, r, "/graph", http.StatusSeeOther)
	})
}

// middleware to check if the user is authenticated
func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if the cookie is set
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}

		provider, err := oidc.NewProvider(r.Context(), "https://accounts.google.com")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
			return
		}

		//verify the token
		verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")})
		_, err = verifier.Verify(r.Context(), cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}

		//if the cookie is set, call the next handler
		next.ServeHTTP(w, r)
	})
}
