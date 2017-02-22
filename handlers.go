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
	H func(redis.Conn, *http.Request) (int, error, interface{})
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := ah.appContext.redisPool.Get()
	defer conn.Close()

	status, err, obj := ah.H(conn, r)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			jsonResponse(w, status, Status{State: "not found"})
		case http.StatusInternalServerError:
			jsonResponse(w, status, Status{State: "internal error"})
		default:
			jsonResponse(w, status, Status{State: http.StatusText(status)})
		}
		return
	}
	jsonResponse(w, status, obj)
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
