package main

import (
	"log"
	"net/http"
	"os"

	"github.com/etix/mirrorbits/database"
)

var redisPool *database.Redis

func main() {
	redisPool = database.NewRedis()

	router := NewRouter()
	var port = ":" + os.Getenv("CDN_API_PORT")
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
