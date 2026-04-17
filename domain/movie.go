package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID
	OrgID     uuid.UUID
	Title     string
	Rating    int
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
