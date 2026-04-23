package db

import (
	"context"

	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

type UsersRepository struct {
	querier generated.Querier
}

func NewUsersRepository(querier generated.Querier) *UsersRepository {
	return &UsersRepository{querier: querier}
}

func (r *UsersRepository) Create(ctx context.Context, user domain.User) error {
	params := generated.CreateUserParams{
		ID:           user.ID.String(),
		OrgID:        user.OrgID.String(),
		Email:        user.Email,
		PasswordHash: user.Password,
		Role:         user.Role.String(),
		Verified:     user.Verified,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	return r.querier.CreateUser(ctx, params)
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return parseUserRow(row)
}

func (r *UsersRepository) Verify(ctx context.Context, userID uuid.UUID) error {
	return r.querier.VerifyUser(ctx, userID.String())
}

func (r *UsersRepository) CreateVerificationToken(ctx context.Context, userID string, token string) error {
	params := generated.CreateVerificationTokenParams{
		UserID: userID,
		Token:  token,
	}
	return r.querier.CreateVerificationToken(ctx, params)
}

func (r *UsersRepository) GetVerificationToken(ctx context.Context, token string) (domain.VerificationToken, error) {
	row, err := r.querier.GetVerificationToken(ctx, token)
	if err != nil {
		return domain.VerificationToken{}, err
	}
	return parseVerificationTokenRow(row)
}

func (r *UsersRepository) DeleteVerificationToken(ctx context.Context, token string) error {
	return r.querier.DeleteVerificationToken(ctx, token)
}

func parseUserRow(row generated.User) (domain.User, error) {
	userID, err := uuid.Parse(row.ID)
	if err != nil {
		return domain.User{}, err
	}

	orgID, err := uuid.Parse(row.OrgID)
	if err != nil {
		return domain.User{}, err
	}

	role, err := domain.ParseRole(row.Role)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:        userID,
		OrgID:     orgID,
		Email:     row.Email,
		Password:  row.PasswordHash,
		Role:      role,
		Verified:  row.Verified,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func parseVerificationTokenRow(row generated.VerificationToken) (domain.VerificationToken, error) {
	userID, err := uuid.Parse(row.UserID)
	if err != nil {
		return domain.VerificationToken{}, err
	}

	return domain.VerificationToken{
		Token:     row.Token,
		UserID:    userID,
		CreatedAt: row.CreatedAt,
	}, nil
}
