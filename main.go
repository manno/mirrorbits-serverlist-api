package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
)

type appContext struct {
	redisPool *redis.Pool
}

var context = &appContext{redisPool: NewPool(":6379")}

func main() {
	router := NewRouter(context)
	url := os.Getenv("API_URL")
	log.Printf("Listening on API_URL from env: %s", url)
	log.Fatal(http.ListenAndServe(url, router))
}
