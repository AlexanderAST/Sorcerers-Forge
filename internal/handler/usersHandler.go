package handler

import (
	"Diploma/configs"
	"Diploma/internal/model"
	"Diploma/internal/validations"
	_ "Diploma/internal/validations"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *server) handleUsersCreate() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, _ := s.store.User().FindByEmail(req.Email)

		if u != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("A user with this email is already registered"))
			return
		}

		configs.RunSmtpRegister(req.Email)

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"status": "success"})
	}
}

func (s *server) handleUsersConfirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &resetPassword{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if ok := checkCode(req.Email, req.EmailCode); !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid emailCode"))
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()

		createProfileEmail := &model.Profile{
			UserID: u.ID,
			Name:   u.Email,
		}

		if err := s.store.Profile().CreateProfile(createProfileEmail); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, map[string]interface{}{"id": u.ID})
	}
}

func checkCode(email, emCode string) bool {
	emailCode := configs.GenToken(email)

	if emailCode == emCode {
		return true
	}
	return false
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := s.GenerateToken(req.Email, req.Password)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{
			"token": token,
		})
	}
}

func (s *server) GenerateToken(email, password string) (string, error) {

	user, err := s.store.User().FindByEmail(email)
	if err != nil {
		return "", err
	}

	if err != nil || !user.ComparePassword(password) {
		return "", errNotAuthenticated
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
	})

	return token.SignedString([]byte(os.Getenv("signing")))
}

func (s *server) ParseToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("signing")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return 0, errors.New("token claims are not of type")
	}
	return claims.UserId, nil
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("userID").(int)

		userEmail, err := s.store.User().Find(userID)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("cannot find this userID: %v", userID)))
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"id": userID, "email": userEmail.Email})
	}
}

func (s *server) sendResetCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &reqWithEmailInput{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		a := validations.ValidateEmail(req.Email)

		if a == false {
			s.error(w, r, http.StatusUnauthorized, errors.New("invalid email format"))
			return
		}

		_, err := s.store.User().SendResetCode(req.Email)

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}

func (s *server) resetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &resetPassword{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		b, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		_, err = s.store.User().ResetPassword(req.Email, req.EmailCode, string(b))

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailCode)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "password successfully changed"})
	}
}

func (s *server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		intID, _ := strconv.Atoi(id)

		_ = s.store.Profile().DeleteProfile(intID)
		err := s.store.User().DeleteUser(intID)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
	}
}
