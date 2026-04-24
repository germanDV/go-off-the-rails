package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/germandv/go-off-the-rails/db"
	"github.com/germandv/go-off-the-rails/db/generated"
	"github.com/germandv/go-off-the-rails/domain"
	"github.com/google/uuid"
)

type AuthController struct {
	mdw         *MiddlewareChain
	dbClient    *sql.DB
	usersRepo   *db.UsersRepository
	orgsRepo    *db.OrgsRepository
	invitesRepo *db.InvitesRepository
	tokenizer   *domain.Tokenizer
}

func NewAuthController(
	mdw *MiddlewareChain,
	dbClient *sql.DB,
	usersRepo *db.UsersRepository,
	orgsRepo *db.OrgsRepository,
	invitesRepo *db.InvitesRepository,
	tokenizer *domain.Tokenizer,
) *AuthController {
	return &AuthController{
		mdw:         mdw,
		dbClient:    dbClient,
		usersRepo:   usersRepo,
		orgsRepo:    orgsRepo,
		invitesRepo: invitesRepo,
		tokenizer:   tokenizer,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /signup", c.mdw.RBACNone(c.Signup))
	mux.Handle("GET /verify", c.mdw.RBACNone(c.Verify))
	mux.Handle("POST /login", c.mdw.RBACNone(c.Login))
	mux.Handle("POST /signout", c.mdw.RBACNone(c.Signout))
	mux.Handle("POST /invites", c.mdw.RBACAdmin(c.CreateInvite))
}

func (c *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	inviteToken := r.FormValue("invite_token")

	passwordHash, err := domain.HashPassword(password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if inviteToken != "" {
		c.signupWithInvite(w, r, email, passwordHash, inviteToken)
		return
	}

	c.signupNewOrg(w, r, email, passwordHash)
}

func (c *AuthController) signupWithInvite(w http.ResponseWriter, r *http.Request, email, passwordHash, inviteToken string) {
	invite, err := c.invitesRepo.GetByToken(r.Context(), inviteToken)
	if err != nil {
		http.Error(w, "invalid invite token", http.StatusBadRequest)
		return
	}

	if invite.IsExpired(time.Now()) {
		http.Error(w, "invite token expired", http.StatusBadRequest)
		return
	}

	if invite.Email != email {
		http.Error(w, "email does not match invite", http.StatusBadRequest)
		return
	}

	userID := uuid.Must(uuid.NewV7())
	user, err := domain.NewUser(userID, invite.OrgID, email, passwordHash, domain.RoleUser)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	user.Verified = true

	tx, err := c.dbClient.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := generated.New(tx)

	err = db.NewUsersRepository(qtx).Create(r.Context(), user)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = db.NewInvitesRepository(qtx).Delete(r.Context(), inviteToken)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Signup successful. You can now log in."))
}

func (c *AuthController) signupNewOrg(w http.ResponseWriter, r *http.Request, email, passwordHash string) {
	userID := uuid.Must(uuid.NewV7())
	orgID := uuid.Must(uuid.NewV7())

	user, err := domain.NewUser(userID, orgID, email, passwordHash, domain.RoleAdmin)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	org, err := domain.NewOrg(orgID, user.Email, time.Now().UTC())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := domain.GenerateVerificationToken(user.ID, time.Now())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	tx, err := c.dbClient.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	qtx := generated.New(tx)

	err = db.NewUsersRepository(qtx).Create(r.Context(), user)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = db.NewOrgsRepository(qtx).Create(r.Context(), org)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = db.NewUsersRepository(qtx).CreateVerificationToken(r.Context(), userID.String(), token.Token)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Signup successful. Please verify your email. Token: %s", token)))
}

func (c *AuthController) CreateInvite(w http.ResponseWriter, r *http.Request) {
	actor := GetActorFromRequest(r)
	if actor == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	inviteID := uuid.Must(uuid.NewV7())
	invite, err := domain.NewInvite(inviteID, actor.OrgID, email, time.Now().UTC())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = c.invitesRepo.Create(r.Context(), invite)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Invite created. Token: %s", invite.Token)))
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.usersRepo.GetByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := domain.CheckPassword(user.Password, password); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if !user.Verified {
		http.Error(w, "need to verify account", http.StatusForbidden)
		return
	}

	actor := domain.Actor{
		UserID: user.ID,
		OrgID:  user.OrgID,
		Role:   user.Role,
		Email:  user.Email,
	}

	jwtString, err := c.tokenizer.Generate(actor)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    jwtString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(domain.TokenExpiration.Seconds()),
	})

	w.Write([]byte("Login successful"))
}

func (c *AuthController) Signout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.Write([]byte("Signed out"))
}

func (c *AuthController) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	vt, err := c.usersRepo.GetVerificationToken(r.Context(), token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	err = c.usersRepo.Verify(r.Context(), vt.UserID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = c.usersRepo.DeleteVerificationToken(r.Context(), token)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
