package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Org struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrg(id uuid.UUID, creatorEmail string, timestamp time.Time) (Org, error) {
	// TODO: validate inputs

	return Org{
		ID:        id,
		Name:      fmt.Sprintf("%s's org", creatorEmail),
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}, nil
}
