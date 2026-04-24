package db

import (
	"context"

	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

type InvitesRepository struct {
	querier generated.Querier
}

func NewInvitesRepository(querier generated.Querier) *InvitesRepository {
	return &InvitesRepository{querier: querier}
}

func (r *InvitesRepository) Create(ctx context.Context, invite domain.Invite) error {
	params := generated.CreateInviteParams{
		ID:        invite.ID.String(),
		OrgID:     invite.OrgID.String(),
		Email:     invite.Email,
		Token:     invite.Token,
		CreatedAt: invite.CreatedAt,
		ExpiresAt: invite.ExpiresAt,
	}
	return r.querier.CreateInvite(ctx, params)
}

func (r *InvitesRepository) GetByToken(ctx context.Context, token string) (domain.Invite, error) {
	row, err := r.querier.GetInviteByToken(ctx, token)
	if err != nil {
		return domain.Invite{}, err
	}

	id, err := uuid.Parse(row.ID)
	if err != nil {
		return domain.Invite{}, err
	}

	orgID, err := uuid.Parse(row.OrgID)
	if err != nil {
		return domain.Invite{}, err
	}

	return domain.Invite{
		ID:        id,
		OrgID:     orgID,
		Email:     row.Email,
		Token:     row.Token,
		CreatedAt: row.CreatedAt,
		ExpiresAt: row.ExpiresAt,
	}, nil
}

func (r *InvitesRepository) Delete(ctx context.Context, token string) error {
	return r.querier.DeleteInvite(ctx, token)
}
