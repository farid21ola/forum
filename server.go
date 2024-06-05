package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/farid21ola/forum/domain"
	"github.com/farid21ola/forum/graph"
	customMiddleware "github.com/farid21ola/forum/middleware"
	"github.com/farid21ola/forum/model"
	"github.com/farid21ola/forum/storage"
	"github.com/farid21ola/forum/storage/postgres"

	"log"
	"net/http"
	"os"
	"time"
)

const (
	localURL    = "postgres://postgres:admin@localhost:5432/forum?sslmode=disable"
	defaultPort = "8080"
)

func main() {
	pgUrl := os.Getenv("DB_URL")
	if pgUrl == "" {
		pgUrl = localURL
	}

	var storage storage.Storage

	pool, err := postgres.NewPoolPostgres(pgUrl)
	if err != nil {
		log.Fatalln("error init DB: ", err)
	}
	storage = postgres.New(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(customMiddleware.AuthMiddleware(storage))

	d := domain.NewDomain(storage)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Domain:            d,
			CommentsObservers: map[string][]chan *model.Comment{},
		}}))

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", graph.DataloaderMiddleware(pool, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
