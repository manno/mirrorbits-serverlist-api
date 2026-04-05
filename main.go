package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
)

type appContext struct {
	redisPool *redis.Pool
}

var context = &appContext{redisPool: NewPool(":6379")}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	router := NewRouter(context)
	url := os.Getenv("API_URL")
	if url == "" {
		url = ":8080"
	}
	_, port, _ := net.SplitHostPort(url)
	slog.Info("starting server", "addr", url, "port", port)
	if err := http.ListenAndServe(url, router); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
