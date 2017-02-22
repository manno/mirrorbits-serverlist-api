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
	// _, err := conn.Do("PING")
	// if err != nil {
	//         panic(err)
	// }
	// exists, err := redis.Bool(conn.Do("EXISTS", "foo"))
	// log.Printf("%s", exists)

	mirrorIDs, err := redis.Strings(conn.Do("LRANGE", "MIRRORS", "0", "-1"))
	if err != nil {
		panic(err)
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
