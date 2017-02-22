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
	router := NewRouter()
	port := ":" + os.Getenv("CDN_API_PORT")
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
