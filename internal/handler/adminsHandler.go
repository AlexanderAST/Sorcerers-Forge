package handler

import (
	"Diploma/internal/model"
	"Diploma/internal/validations"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *server) handleAdminsCreate() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requestAdmin{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.Admin().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": u.ID})
	}
}

func (s *server) handleSessionsAdminsCreate() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requestAdmin{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.Admin().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}

func (s *server) authenticateAdmins(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.Admin().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleAdminsWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) sendResetCodeAdmins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &reqWithEmail{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		a := validations.ValidateEmail(req.Email)

		if a == false {
			s.error(w, r, http.StatusUnauthorized, errors.New("invalid email format"))
			return
		}

		_, err := s.store.Admin().SendResetCode(req.Email)

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}

func (s *server) resetAdminsPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &resetAdminPassword{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		b, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err = s.store.Admin().ResetPassword(req.Email, req.EmailCode, string(b))

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "password successfully changed"})
	}
}

func (s *server) deleteAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &reqWithEmail{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err := s.store.Admin().DeleteUser(req.Email)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}
