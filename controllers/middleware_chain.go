package controllers

import (
	"net/http"
	"slices"

	"github.com/germandv/go-off-the-rails/domain"
)

type Chain []func(http.Handler) http.Handler

func (c Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

func (c Chain) then(h http.Handler) http.Handler {
	for _, mdw := range slices.Backward(c) {
		h = mdw(h)
	}
	return h
}

type MiddlewareChain struct {
	anonChain  Chain
	userChain  Chain
	adminChain Chain
	godChain   Chain
	RBACNone   func(http.HandlerFunc) http.Handler
	RBACUser   func(http.HandlerFunc) http.Handler
	RBACAdmin  func(http.HandlerFunc) http.Handler
	RBACGod    func(http.HandlerFunc) http.Handler
}

func NewMiddlewareChain(tokenizer *domain.Tokenizer) *MiddlewareChain {
	anonChain := Chain{
		Recover,

		// TODO: add these middlewares
		// RealIP,
		// Helmet,
		// COP,
		// ParseForm,

		DetectAuth(tokenizer),
	}

	userChain := append(anonChain, RBAC([]domain.Role{domain.RoleUser}))
	adminChain := append(anonChain, RBAC([]domain.Role{domain.RoleUser, domain.RoleAdmin}))
	godChain := append(anonChain, RBAC([]domain.Role{domain.RoleUser, domain.RoleAdmin, domain.RoleGod}))

	return &MiddlewareChain{
		RBACNone:  anonChain.ThenFunc,
		RBACUser:  userChain.ThenFunc,
		RBACAdmin: adminChain.ThenFunc,
		RBACGod:   godChain.ThenFunc,
	}
}
