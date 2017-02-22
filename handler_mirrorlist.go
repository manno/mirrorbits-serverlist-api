package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/manno/mirrorbits-api/mirror"
	"net/http"
)

type MirrorListResponse struct {
	Status     Status
	MirrorList []mirror.Mirror
}

func Index(conn redis.Conn, r *http.Request) (int, error, interface{}) {
	mirrorIDs, err := redis.Strings(conn.Do("LRANGE", "MIRRORS", "0", "-1"))
	if err != nil {
		return http.StatusInternalServerError, err, nil
	}
	fmt.Printf("%v\n", mirrorIDs)

	mirrors := make([]mirror.Mirror, 0, len(mirrorIDs))
	for _, id := range mirrorIDs {
		mirror := mirror.Mirror{ID: id}
		mirrors = append(mirrors, mirror)
	}

	m := MirrorListResponse{Status: Status{State: OkStatus}, MirrorList: mirrors}
	return http.StatusOK, err, m
}
