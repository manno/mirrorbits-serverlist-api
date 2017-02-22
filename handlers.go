package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"net/http"
)

const OkStatus string = "ok"

type Status struct {
	State string `json:"status"`
}

type appHandler struct {
	*appContext
	H func(redis.Conn, http.ResponseWriter, *http.Request)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := ah.appContext.redisPool.Get()
	defer conn.Close()

	ah.H(conn, w, r)
}

func jsonResponse(w http.ResponseWriter, code int, obj interface{}) {
	writeJSONHeader(w, code)
	jsonEncode(w, obj)
}

func writeJSONHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
}

func jsonEncode(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		panic(err)
	}
}
