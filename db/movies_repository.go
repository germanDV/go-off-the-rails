package db

import (
	"context"

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

func (r *MoviesRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Movie, error) {
	row, err := r.querier.GetMovie(ctx, id.String())
	if err != nil {
		return domain.Movie{}, err
	}
	return parseRow(row)
}

func (r *MoviesRepository) List(ctx context.Context, orgID uuid.UUID, pagination domain.PaginationParams) ([]domain.Movie, error) {
	rows, err := r.querier.ListMovies(ctx, orgID.String(), pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, err
	}
	return mapRows(rows, parseRow)
}

func parseRow(row generated.MovieRow) (domain.Movie, error) {
	return domain.Movie{
		ID:        row.ID,
		OrgID:     row.OrgID,
		Title:     row.Title,
		Rating:    int(row.Rating),
		Version:   int(row.Version),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
