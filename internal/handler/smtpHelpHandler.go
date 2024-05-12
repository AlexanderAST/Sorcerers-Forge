package handler

import (
	configss "Diploma/configs"
	"Diploma/internal/model"
	"encoding/json"
	"net/http"
)

func (s *server) sendHelp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := &Message{}

		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		c := &model.Message{
			UserEmail: msg.UserEmail,
			Name:      msg.Name,
			Message:   msg.Message,
		}

		if _, err := configss.RunSmtpHelp(c); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})

	}
}
