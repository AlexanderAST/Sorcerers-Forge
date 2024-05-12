package handler

import (
	"Diploma/internal/model"
	"encoding/json"
	"errors"
	"net/http"
)

func (s *server) createGallery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &gallery{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if input.Photo == "" || input.Description == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid data"))
			return
		}

		in := &model.Gallery{
			Photo:       input.Photo,
			Description: input.Description,
		}

		if err := s.store.Gallery().CreateGallery(in); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "success"})

	}
}

func (s *server) deleteGallery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &gallery{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Gallery().DeleteGallery(input.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"delete": "success"})
	}
}

func (s *server) getAllGallery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gallerys, err := s.store.Gallery().GetAllGallery()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"gallery": gallerys})
	}
}

func (s *server) updateGallery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &gallery{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		in := &model.Gallery{
			ID:          input.ID,
			Photo:       input.Photo,
			Description: input.Description,
		}

		if _, err := s.store.Gallery().UpdateGallery(in); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}
