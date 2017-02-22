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
	port := ":" + os.Getenv("API_PORT")
	log.Printf("Listening on API_PORT from env: %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
