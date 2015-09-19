package main

import (
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"net/http"
)

func rateLimit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn := pool.Get()
		defer conn.Close()

		remoteIP := r.Header.Get("REMOTE_ADDR")
		path := r.URL.Path // Likely to be more generic to stop attempts across multiple users

		// Increment counter for request (this will create a new key if one does not exist)
		current, err := redis.Int(conn.Do("INCR", path+remoteIP))
		if err != nil {
			LogError(err.Error())
		}

		// Check if the returned counter exceeds our limit
		if int(current) > int(ConRateLimitLimit) {
			w.WriteHeader(429)
		} else if current == 1 {
			// Set the expiry on a fresh counter for the given path and remote address
			conn.Do("EXPIRE", path+remoteIP, ConRateLimitTimeout)
			if err != nil {
				LogError(err.Error())
			}
			h.ServeHTTP(w, r)
		} else {
			// Continue to the next handler if we haven't exceeded our limit
			h.ServeHTTP(w, r)
		}
	})
}
