package handler

import (
	"Diploma/configs"
	"net/http"
)

func (s *server) sendUserOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		carts, err := s.store.Catalog().GetAllCartUser(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = s.store.Catalog().CreateOrder(carts)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		email, _ := s.store.User().Find(userID)

		if _, err := configs.RunSmtpOrders(email.Email, carts); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.Catalog().DeleteAllFromCart(userID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"static": "successfully send orders"})

	}
}

func (s *server) createOrderHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		products, err := s.store.Catalog().GetUserOrderHistory(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"orders": products})
	}
}
