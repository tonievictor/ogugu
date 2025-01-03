package main

import (
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/tonievictor/dotenv"

	"ogugu/docs"
	"ogugu/router"
)

func main() {
	dotenv.Config()

	docs.SwaggerInfo.Title = "Ogugu"
	docs.SwaggerInfo.Description = "An RSS feed reader"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "localhost:8080" // this will be dynamic
	docs.SwaggerInfo.BasePath = "/v1/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := router.Routes()

	http.ListenAndServe(os.Getenv("PORT"), r)
}
