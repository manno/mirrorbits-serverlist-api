package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Mirror struct {
	ID             string
	Name           string
	HttpURL        string  `redis:"http"`
	Latitude       float32 `redis:"latitude"`
	Longitude      float32 `redis:"longitude"`
	ContinentCode  string  `redis:"continentCode"`
	CountryCodes   string  `redis:"countryCodes"`
	LastSync       int64   `redis:"lastSync"`
	Asnum          int     `redis:"asnum"`
	SponsorURL     string  `redis:"sponsorURL"`
	SponsorLogoURL string  `redis:"sponsorLogo"`
	SponsorName    string  `redis:"sponsorName"`
	Up             int     `redis:"up"`
	Enabled        bool    `redis:"enabled"`
	FileCount      int64
	MonthDownloads int64
	MonthBytes     int64
}

type MirrorListResponse struct {
	Status     Status
	MirrorList []Mirror
}

type MirrorsSlice []Mirror

func (s MirrorsSlice) Len() int      { return len(s) }
func (s MirrorsSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByDownloadNumbers struct {
	MirrorsSlice
}

func (b ByDownloadNumbers) Less(i, j int) bool {
	if b.MirrorsSlice[i].MonthDownloads > b.MirrorsSlice[j].MonthDownloads {
		return true
	}
	return false
}

func Index(conn redis.Conn, r *http.Request) (int, error, interface{}) {
	mlist, err := redis.Strings(conn.Do("HGETALL", "MIRRORS"))
	if err != nil {
		return http.StatusInternalServerError, err, nil
	}

	mirrors := make([]Mirror, 0, len(mlist)/2)
	for i := 0; i < len(mlist)-1; i = i + 2 {
		mirrors = append(mirrors, Mirror{ID: mlist[i], Name: mlist[i+1]})
	}

	conn.Send("MULTI")
	for _, m := range mirrors {
		month := time.Now().Format("2006_01")
		conn.Send("HGET", "STATS_MIRROR_"+month, m.ID)
		conn.Send("HGET", "STATS_MIRROR_BYTES_"+month, m.ID)
	}
	stats, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		return http.StatusInternalServerError, err, nil
	}

	enabled := make([]Mirror, 0, len(mlist)/2)
	index := 0
	for _, mirror := range mirrors {
		downloads, _ := redis.Int64(stats[index], nil)
		bytes, _ := redis.Int64(stats[index+1], nil)
		index += 2

		count, _ := redis.Int64(conn.Do("SCARD", fmt.Sprintf("MIRRORFILES_%s", mirror.ID)))

		mirror.MonthDownloads = downloads
		mirror.MonthBytes = bytes
		mirror.FileCount = count

		reply, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("MIRROR_%s", mirror.ID)))
		if err != nil {
			continue
		}
		if len(reply) == 0 {
			err = redis.ErrNil
			continue
		}
		err = redis.ScanStruct(reply, &mirror)
		if err != nil {
			continue
		}
		if !mirror.Enabled {
			continue
		}

		mirror.ID = mirror.Name
		enabled = append(enabled, mirror)
	}

	sort.Sort(ByDownloadNumbers{mirrors})
	m := MirrorListResponse{Status: Status{State: OkStatus}, MirrorList: enabled}
	return http.StatusOK, err, m
}
