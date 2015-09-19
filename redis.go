package main

import (
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
)

var pool redis.Pool

func createRedisPool() {
	pool = redis.Pool{
		MaxIdle:   ConMaxRedisIdle,
		MaxActive: ConMaxRedisActive, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ConRedisUrl)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
