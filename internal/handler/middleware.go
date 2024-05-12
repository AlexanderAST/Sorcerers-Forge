package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func (s *server) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			s.error(w, r, http.StatusUnauthorized, errors.New("empty auth header"))
			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 {
			s.error(w, r, http.StatusUnauthorized, errors.New("invalid auth header"))
			return
		}
		userId, err := s.ParseToken(headerParts[1])
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *server) optionalUserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		userId, err := s.ParseToken(headerParts[1])
		if err != nil {
			// Пропускаем запрос с неверным токеном
			next.ServeHTTP(w, r)
			return
		}

		// Устанавливаем userID в контексте
		ctx := context.WithValue(r.Context(), "userID", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
