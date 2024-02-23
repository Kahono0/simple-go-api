package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Kahono0/simple-go-api/graph"
	"github.com/Kahono0/simple-go-api/internal"
	"github.com/Kahono0/simple-go-api/utils"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const defaultPort = "8080"

func indexHandler(w http.ResponseWriter, r *http.Request) {

	loginHtml := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Login</title>
		</head>
		<body>
			<h1>Login</h1>
			<form action="/auth" method="post">
				<input type="submit" value="Login with Google">
			</form>
		</body>
		</html>
	`

	//check if token is present
	_, err := r.Cookie("token")
	if err != nil {
		w.Write([]byte(loginHtml))
		return
	}

	http.Redirect(w, r, "/graph", http.StatusTemporaryRedirect)
}

func main() {

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//setup database
	utils.ConnectDB()

	ctx := context.Background()

	//openid connect config
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH2_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	//setup routes
	http.HandleFunc("/", indexHandler)

	//health check
	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	//setup authentication
	internal.SetUpAuth(&config, provider)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/graph", internal.AuthMiddleWare(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", internal.AuthMiddleWare(srv))

	log.Printf("connect to http://127.0.0.1:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
