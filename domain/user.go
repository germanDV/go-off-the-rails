package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	OrgID     uuid.UUID
	Email     string
	Password  string
	Role      Role
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id uuid.UUID, orgID uuid.UUID, email, password string, role Role) (User, error) {
	// TODO: validate inputs

	return User{
		ID:        id,
		OrgID:     orgID,
		Email:     email,
		Password:  password,
		Role:      role,
		Verified:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) IsVerified() bool {
	return u.Verified
}

func (u *User) Verify(timestamp time.Time) {
	u.Verified = true
	u.UpdatedAt = timestamp
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
