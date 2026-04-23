package controllers

import (
	"context"
	"net/http"
	"slices"

	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

const (
	AuthCookieName = "auth_token"
)

type CtxKey string

const UserKey CtxKey = "user"

// DetectAuth reads the auth_token cookie.
// If it is present, validates the JWT, and sets user data in the request context.
// Otherwise, it calls next without any modifications to the request.
func DetectAuth(tokenizer *domain.Tokenizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(AuthCookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := tokenizer.Validate(cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := uuid.Parse(claims["sub"].(string))
			if err != nil || userID == uuid.Nil {
				next.ServeHTTP(w, r)
				return
			}

			orgID, err := uuid.Parse(claims["org_id"].(string))
			if err != nil || orgID == uuid.Nil {
				next.ServeHTTP(w, r)
				return
			}

			role, err := domain.ParseRole(claims["role"].(string))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			email, _ := claims["email"].(string)

			actor := domain.Actor{
				UserID: userID,
				OrgID:  orgID,
				Role:   role,
				Email:  email,
			}

			r = r.WithContext(context.WithValue(r.Context(), UserKey, &actor))
			next.ServeHTTP(w, r)
		})
	}
}

// RBAC checks that a user is present in the request context and has one of the given roles.
// If the user is not present or does not have one of the given roles, it redirects to the login page.
func RBAC(roles []domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor := GetActorFromRequest(r)
			if actor == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if !slices.Contains(roles, actor.Role) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetActorFromRequest returns the Actor from the request context.
// If the Actor is not set, it returns nil.
func GetActorFromRequest(r *http.Request) *domain.Actor {
	a, ok := r.Context().Value(UserKey).(*domain.Actor)
	if !ok {
		return nil
	}
	if a == nil || a.UserID == uuid.Nil || a.OrgID == uuid.Nil {
		return nil
	}
	return a
}
