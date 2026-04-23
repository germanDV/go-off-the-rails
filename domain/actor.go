package domain

import "github.com/google/uuid"

type Actor struct {
	UserID uuid.UUID
	OrgID  uuid.UUID
	Role   Role
	Email  string
}
