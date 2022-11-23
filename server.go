package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/mondracode/ambrosia-atlas-api/graph"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/rs/cors"
)

const defaultPort = "5000"

func main() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Use(middleware.SetHeader("Origin", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/graphql", srv)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}
