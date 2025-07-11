package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

type key string

const authKey key = "authID"

type Middleware func(http.Handler) http.Handler

func addMiddleware(h http.Handler) http.Handler {
	middlewares := []Middleware{
		requestLoggerMiddleware,
		authMiddleware,
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

type authInfo struct {
	authenticated bool
	userID        uint64
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("x-auth-id")
		info := authInfo{authenticated: false, userID: 0}
		if authHeader != "" {
			id, err := strconv.Atoi(authHeader)
			if err == nil {
				info.authenticated = true
				info.userID = uint64(id) // TODO: need to actually verify this
			} else {
				log.Panicf("failed to parse x-auth-id %v", err)
			}
		}
		ctx := context.WithValue(r.Context(), authKey, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("req = %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
