package main

import (
	"log/slog"
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
	slog.Info("starting server", "addr", url)
	if err := http.ListenAndServe(url, router); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
