package handler

import (
	"Diploma/internal/store"
	"Diploma/internal/store/sqlstore"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

const (
	sessionName        = "cookie"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func Start() error {
	db, err := newPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslMode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("session_key")))
	srv := newServer(store, sessionStore)

	logrus.Info("Server start")

	portConfig := viper.GetString("port")

	return http.ListenAndServe(portConfig, srv)
}

func newPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
