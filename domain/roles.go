package domain

import "fmt"

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
	RoleGod   Role = "GOD"
)

func (s Role) String() string {
	return string(s)
}

func ParseRole(s string) (Role, error) {
	switch s {
	case RoleUser.String():
		return RoleUser, nil
	case RoleAdmin.String():
		return RoleAdmin, nil
	case RoleGod.String():
		return RoleGod, nil
	default:
		return "", fmt.Errorf("invalid role: %s", s)
	}
}
