package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/germandv/go-off-the-rails/db"
	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
)

const (
	HttpIdleTimeout       = 1 * time.Minute
	HttpReadHeaderTimeout = 5 * time.Second
	HttpReadTimeout       = 10 * time.Second
	HttpWriteTimeout      = 15 * time.Second
	HttpTimeout           = 14 * time.Second
)

type ServerConfig struct {
	Port      int
	DBClient  *sql.DB
	Tokenizer *domain.Tokenizer
}

type Server struct {
	dbClient   *sql.DB
	mdw        *MiddlewareChain
	httpServer *http.Server
	mux        *http.ServeMux
}

func NewServer(cfg ServerConfig) *Server {
	mdw := NewMiddlewareChain(cfg.Tokenizer)
	mux := &http.ServeMux{}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		IdleTimeout:       HttpIdleTimeout,
		ReadHeaderTimeout: HttpReadHeaderTimeout,
		ReadTimeout:       HttpReadTimeout,
		WriteTimeout:      HttpWriteTimeout,
	}

	httpServer.Handler = http.TimeoutHandler(mux, HttpTimeout, "Request timed out")

	s := &Server{
		dbClient:   cfg.DBClient,
		mdw:        mdw,
		mux:        mux,
		httpServer: httpServer,
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	healthController := NewHealthController(s.mdw)
	healthController.RegisterRoutes(s.mux)

	moviesRepo := db.NewMoviesRepository(generated.New(s.dbClient))
	moviesController := NewMoviesController(s.mdw, moviesRepo)
	moviesController.RegisterRoutes(s.mux)
}

func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
