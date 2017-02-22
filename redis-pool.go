package main

import "github.com/garyburd/redigo/redis"
import "time"

// NewPool as described in https://godoc.org/github.com/garyburd/redigo/redis#Pool
func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}
