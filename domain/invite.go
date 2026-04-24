package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const InviteExpiration = 7 * 24 * time.Hour

type Invite struct {
	ID        uuid.UUID
	OrgID     uuid.UUID
	Email     string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func GenerateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewInvite(id, orgID uuid.UUID, email string, now time.Time) (Invite, error) {
	token, err := GenerateInviteToken()
	if err != nil {
		return Invite{}, err
	}
	return Invite{
		ID:        id,
		OrgID:     orgID,
		Email:     email,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now.Add(InviteExpiration),
	}, nil
}

func (i *Invite) IsExpired(now time.Time) bool {
	return now.After(i.ExpiresAt)
}
