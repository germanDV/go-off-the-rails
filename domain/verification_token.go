package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type VerificationToken struct {
	Token     string
	UserID    uuid.UUID
	CreatedAt time.Time
}

func GenerateVerificationToken(userID uuid.UUID, timestamp time.Time) (VerificationToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return VerificationToken{}, err
	}

	return VerificationToken{
		Token:     hex.EncodeToString(b),
		UserID:    userID,
		CreatedAt: timestamp,
	}, nil
}
