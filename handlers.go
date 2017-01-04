package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const OkStatus string = "ok"
const JsonErrorStatus string = "json invalid"
const InvalidKeyStatus string = "api key invalid"

type Status struct {
	State string `json:"status"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	conn, err := redisPool.Connect()
	if err != nil {
		log.Fatal("Redis: ", err)
	}
	defer conn.Close()

	jsonResponse(w, http.StatusOK, OkStatus)
}

func jsonEncode(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		panic(err)
	}
}

func jsonResponse(w http.ResponseWriter, code int, status string) {
	writeJsonHeader(w, code)
	msg := Status{State: status}
	jsonEncode(w, msg)
}

func writeJsonHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
}
