package handler

import (
	"Diploma/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (s *server) createProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &profileInput{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		userID := r.Context().Value("userID").(int)

		if input.Name == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid data"))
			return
		}

		p := &model.Profile{
			UserID:     userID,
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
			Contact:    input.Contact,
			Photo:      input.Photo,
		}

		if err := s.store.Profile().CreateProfile(p); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": p.ID, "status": "success"})
	}
}

func (s *server) updateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &profileInput{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		p := &model.Profile{
			UserID:     userID,
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
			Contact:    input.Contact,
			Photo:      input.Photo,
		}

		if _, err := s.store.Profile().UpdateProfile(p); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": userID, "status": "success"})
	}
}

func (s *server) deleteProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		intId, _ := strconv.Atoi(id)

		if err := s.store.Profile().DeleteProfile(intId); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}

func (s *server) takeProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("userID").(int)

		profile, err := s.store.Profile().FindByID(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"profile": profile})
	}
}
