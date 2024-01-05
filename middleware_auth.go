package main

import (
	"net/http"
	"strings"

	"github.com/tyler-lutz/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		vals := strings.Split(authHeader, " ")
		if len(vals) != 2 {
			respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header")
			return
		}
		apiKey := vals[1]

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
			return
		}

		handler(w, r, user)
	}
}
