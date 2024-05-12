package handler

import (
	"Diploma/configs"
	"Diploma/internal/model"
	"encoding/json"
	"net/http"
)

func (s *server) createLearnRequestNoAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.Message{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userReq := &model.Message{
			UserEmail: req.UserEmail,
			Name:      req.Name,
			Message:   req.Message,
		}

		_, err := configs.RunSmtpLearn(userReq)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully send request"})
	}
}

func (s *server) createLearnRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.Message{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		userID := r.Context().Value("userID").(int)

		email, err := s.store.User().Find(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		name, err := s.store.Profile().FindByID(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userReq := &model.Message{
			UserEmail: email.Email,
			Name:      name.Name,
			Message:   req.Message,
		}

		_, err = configs.RunSmtpLearn(userReq)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "successfully send request"})
	}
}

func (s *server) learnWrapper(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" {
		s.createLearnRequestNoAuth()(w, r)
		return
	}

	s.createLearnRequest()(w, r)
}
