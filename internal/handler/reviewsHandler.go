package handler

import (
	"Diploma/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (s *server) createReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &reviews{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		hasOrder, err := s.store.Catalog().GetUserOrderHistoryByReviews(userID, input.ProductID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		fmt.Println(hasOrder)
		if !hasOrder {
			s.error(w, r, http.StatusBadRequest, errors.New("you can't send review"))
			return
		}

		p := &model.Reviews{
			ProductId: input.ProductID,
			UserID:    userID,
			Stars:     input.Stars,
			Message:   input.Message,
		}

		if err := s.store.Reviews().CreateReviews(p); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": p.ID, "status": "success"})
	}
}

func (s *server) getAllReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reviews, err := s.store.Reviews().GetAllReviews()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"reviews": reviews})
	}
}

func (s *server) getAllReviewsByProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("id")
		productID, _ := strconv.Atoi(product)

		review, err := s.store.Reviews().GetAllReviewsFromProduct(productID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"reviews": review})
	}
}

func (s *server) updateReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &reviews{}

		if err := json.NewDecoder(r.Body).Decode(input); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		userID := r.Context().Value("userID").(int)

		p := &model.Reviews{
			ProductId: input.ProductID,
			UserID:    userID,
			Stars:     input.Stars,
			Message:   input.Message,
		}

		if _, err := s.store.Reviews().UpdateReviews(p); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "success"})
	}
}

func (s *server) deleteReview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("userID").(int)
		id := r.URL.Query().Get("id")
		productID, _ := strconv.Atoi(id)

		if err := s.store.Reviews().DeleteReview(productID, userID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "delete success"})
	}
}
