package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Kahono0/simple-go-api/graph"
)

const defaultPort = "8080"

func authenticate(w http.ResponseWriter, r *http.Request) {
	//set token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: "token",
	})
	w.WriteHeader(http.StatusOK)
	//redirect to the graphiql playground
	http.Redirect(w, r, "/graph", http.StatusSeeOther)
}

func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if the cookie is set
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		//check if the cookie value is correct
		if cookie.Value != "token" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		//if the cookie is set and the value is correct, call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	http.Handle("/auth", http.HandlerFunc(authenticate))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/graph", authMiddleWare(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", authMiddleWare(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
