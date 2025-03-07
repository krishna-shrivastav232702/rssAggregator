package main

import (
	"log"
	"database/sql"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/krishna-shrivastav232702/rssAggregator/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL is not found in the environment")
	}

	conn,err := sql.Open("postgres",dbURL)
	if err != nil {
		log.Fatal("Cant connect to database")
	}
	apiCfg := apiConfig{
		DB:database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*","http://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)	
	v1Router.Get("/error",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds",apiCfg.handlerGetFeeds)
	v1Router.Post("/feedFollows",apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feedFollows",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feedFollows/{feedFollowID}",apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	router.Mount("/v1",v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
