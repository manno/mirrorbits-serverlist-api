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

type MirrorList struct {
	Status     Status
	MirrorList []Mirror
}

type Mirror struct {
	ID             string
	HttpURL        string
	Latitude       float32
	Longitude      float32
	ContinentCode  string
	CountryCodes   string
	LastSync       int64
	Asnum          int
	SponsorURL     string
	SponsorLogoURL string
	SponsorName    string
	FileCount      int
}

func Index(w http.ResponseWriter, r *http.Request) {
	conn, err := redisPool.Connect()
	if err != nil {
		log.Fatal("Redis: ", err)
	}
	defer conn.Close()

	m := MirrorList{Status: Status{State: OkStatus}}
	jsonResponse(w, http.StatusOK, m)
}

func jsonEncode(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		panic(err)
	}
}

func jsonResponse(w http.ResponseWriter, code int, msg MirrorList) {
	writeJsonHeader(w, code)
	jsonEncode(w, msg)
}

func writeJsonHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
}
