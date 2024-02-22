package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/Kahono0/simple-go-api/engine"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Claims struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
}

func verifyToken(token string) (*Claims, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")})
	v, err := verifier.Verify(context.Background(), token)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	err = v.Claims(claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

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

		userInfo, err := provider.UserInfo(r.Context(), oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(w, "failed to get user info", http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "no id_token", http.StatusInternalServerError)
			return
		}

		claims, err := verifyToken(rawIDToken)
		if err != nil {
			http.Error(w, "failed to verify token", http.StatusInternalServerError)
			return
		}

		fmt.Println(claims.FamilyName)

		//create user if not exists
		err = engine.CreateUser(claims.Sub, claims.Email, claims.Name)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		//set the token in the cookie
		c := &http.Cookie{
			Name:    "token",
			Value:   rawIDToken,
			MaxAge:  24 * 60 * 60,
			Path:    "/",
			Expires: token.Expiry,
		}

		http.SetCookie(w, c)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, %v!", userInfo.Subject)

	})
}

// middleware to check if the user is authenticated
func AuthMiddleWare(next http.Handler) http.Handler {
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
		v, err := verifier.Verify(r.Context(), cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized!")
			return
		}

		claims := &Claims{}
		err = v.Claims(claims)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized!")
			return
		}

		user, err := engine.GetUserByID(claims.Sub)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
			return
		}

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized!")
			return
		}

		//set the user in the context
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		//if the cookie is set, call the next handler
		next.ServeHTTP(w, r)
	})
}
