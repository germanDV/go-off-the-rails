package db

import (
	"context"
	"database/sql"

	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

type MoviesRepository struct {
	querier generated.Querier
}

func NewMoviesRepository(querier generated.Querier) *MoviesRepository {
	return &MoviesRepository{querier: querier}
}

func (r *MoviesRepository) Create(ctx context.Context, movie domain.Movie) error {
	params := generated.CreateMovieParams{
		ID:    movie.ID.String(),
		OrgID: movie.OrgID.String(),
		Title: movie.Title,
		Rating: sql.NullInt64{
			Int64: movie.Rating,
			Valid: movie.Rating != 0,
		},
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}
	return r.querier.CreateMovie(ctx, params)
}

func (r *MoviesRepository) GetByID(
	ctx context.Context,
	orgID uuid.UUID,
	movieID uuid.UUID,
) (domain.Movie, error) {
	params := generated.GetMovieParams{
		ID:    movieID.String(),
		OrgID: orgID.String(),
	}

	row, err := r.querier.GetMovie(ctx, params)
	if err != nil {
		return domain.Movie{}, err
	}

	return parseRow(row)
}

func (r *MoviesRepository) List(
	ctx context.Context,
	orgID uuid.UUID,
	pagination domain.PaginationParams,
) ([]domain.Movie, error) {
	params := generated.ListMoviesParams{
		OrgID:  orgID.String(),
		Limit:  pagination.Limit(),
		Offset: pagination.Offset(),
	}

	rows, err := r.querier.ListMovies(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapRows(rows, parseRow)
}

func parseRow(row generated.Movie) (domain.Movie, error) {
	movieID, err := uuid.Parse(row.ID)
	if err != nil {
		return domain.Movie{}, err
	}

	orgID, err := uuid.Parse(row.OrgID)
	if err != nil {
		return domain.Movie{}, err
	}

	return domain.Movie{
		ID:        movieID,
		OrgID:     orgID,
		Title:     row.Title,
		Rating:    row.Rating.Int64,
		Version:   row.Version,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
