package handler

import (
	"Diploma/internal/model"
	"encoding/json"
	"errors"
	"net/http"
)

func (s *server) createAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &address{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if input.Name == "" || input.Latlng == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid name or latlng"))
			return
		}

		c := &model.Address{
			Name:   input.Name,
			Latlng: input.Latlng,
		}

		if err := s.store.Address().CreateAddress(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": c.ID})

	}

}

func (s *server) getAllAddresses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addresses, err := s.store.Address().GetAllAddresses()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"addresses": addresses})
	}
}

func (s *server) deleteAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &address{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Address().DeleteAddress(input.ID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"delete": "success"})
	}
}

func (s *server) updateAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &address{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		c := &model.Address{
			ID:     input.ID,
			Name:   input.Name,
			Latlng: input.Latlng,
		}
		if _, err := s.store.Address().UpdateAddress(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}
