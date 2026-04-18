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
	mux.HandleFunc("POST /movies", c.Create)
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
	movieID, err := uuid.Parse(r.PathValue("movie_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: get form user in r.Context()
	orgID := uuid.Must(uuid.NewV7())

	movie, err := c.moviesRepo.GetByID(r.Context(), orgID, movieID)
	if err != nil {
		// TODO: handle not found
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("We found movie %s in your collection", movie.Title)))
}

func (c *MoviesController) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: get form user in r.Context()
	orgID := uuid.Must(uuid.NewV7())

	title := r.FormValue("title")
	rating := r.FormValue("rating")

	movie, err := domain.NewMovie(orgID, title, rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.moviesRepo.Create(r.Context(), movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("We created movie %s in your collection", movie.ID)))
}
