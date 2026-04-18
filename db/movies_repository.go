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

func (r *MoviesRepository) Create(ctx context.Context, movie domain.Movie) error {
	row := generated.MovieRow{
		ID:        movie.ID.String(),
		OrgID:     movie.OrgID.String(),
		Title:     movie.Title,
		Rating:    movie.Rating,
		Version:   1,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}
	return r.querier.CreateMovie(ctx, row)
}

func (r *MoviesRepository) GetByID(
	ctx context.Context,
	orgID uuid.UUID,
	movieID uuid.UUID,
) (domain.Movie, error) {
	row, err := r.querier.GetMovie(ctx, orgID.String(), movieID.String())
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
	rows, err := r.querier.ListMovies(ctx, orgID.String(), pagination.Limit(), pagination.Offset())
	if err != nil {
		return nil, err
	}
	return mapRows(rows, parseRow)
}

func parseRow(row generated.MovieRow) (domain.Movie, error) {
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
		Rating:    int(row.Rating),
		Version:   int(row.Version),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
