package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/germandv/go-off-the-rails/controllers"
	"github.com/germandv/go-off-the-rails/db"
	"github.com/germandv/go-off-the-rails/db/generated"
)

const (
	HttpPort              = 8787
	HttpIdleTimeout       = 1 * time.Minute
	HttpReadHeaderTimeout = 5 * time.Second
	HttpReadTimeout       = 10 * time.Second
	HttpWriteTimeout      = 15 * time.Second
	HttpTimeout           = 14 * time.Second
)

func main() {
	moviesRepo := db.NewMoviesRepository(generated.MockQuerier{})
	moviesController := controllers.NewMoviesController(moviesRepo)

	mux := &http.ServeMux{}
	moviesController.RegisterRoutes(mux)

	port := HttpPort
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		IdleTimeout:       HttpIdleTimeout,
		ReadHeaderTimeout: HttpReadHeaderTimeout,
		ReadTimeout:       HttpReadTimeout,
		WriteTimeout:      HttpWriteTimeout,
	}

	httpServer.Handler = http.TimeoutHandler(mux, HttpTimeout, "Request timed out")

	fmt.Printf("Listening on port %d\n", port)
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
