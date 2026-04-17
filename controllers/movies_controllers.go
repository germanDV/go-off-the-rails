package controllers

import (
	"fmt"
	"net/http"

	"github.com/germandv/go-off-the-rails/db"
	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

type MoviesController struct {
	moviesRepo *db.MoviesRepository
}

func NewMoviesController(moviesRepo *db.MoviesRepository) *MoviesController {
	return &MoviesController{moviesRepo: moviesRepo}
}

func (c *MoviesController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /movies", c.Index)
	mux.HandleFunc("GET /movies/{movie_id}", c.Show)
}

func (c *MoviesController) Index(w http.ResponseWriter, r *http.Request) {
	page := domain.NewPaginationParams(
		r.URL.Query().Get("limit"),
		r.URL.Query().Get("offset"),
	)

	// TODO: get form user in r.Context()
	orgID := uuid.Must(uuid.NewV7())

	movies, err := c.moviesRepo.List(r.Context(), orgID, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("We found %d movies in your collection", len(movies))))
}

func (c *MoviesController) Show(w http.ResponseWriter, r *http.Request) {
	movieIDStr := r.PathValue("movie_id")
	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: get form user in r.Context()
	orgID := uuid.Must(uuid.NewV7())

	movie, err := c.moviesRepo.GetByID(r.Context(), movieID)
	if err != nil {
		// TODO: handle not found
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if movie.OrgID != orgID {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(fmt.Sprintf("We found movie %s in your collection", movie.Title)))
}
