package db

import (
	"context"

	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
)

type OrgsRepository struct {
	querier generated.Querier
}

func NewOrgsRepository(querier generated.Querier) *OrgsRepository {
	return &OrgsRepository{querier: querier}
}

func (r *OrgsRepository) Create(ctx context.Context, org domain.Org) error {
	params := generated.CreateOrgParams{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
	return r.querier.CreateOrg(ctx, params)
}
