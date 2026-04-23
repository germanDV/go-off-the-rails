package domain

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenIssuer     = "https://gotr.dev"
	TokenAudience   = "https://gotr.dev"
	TokenExpiration = time.Hour * 24 * 7
)

// Tokenizer provides functionality for generating and validating JWTs using ECDSA private and public keys.
type Tokenizer struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewTokenizer initializes and returns a Tokenizer object using the given private and public keys in PEM format.
func NewTokenizer(priv string, publ string) (*Tokenizer, error) {
	privateKey, err := decodePrivate(priv)
	if err != nil {
		return nil, err
	}

	publicKey, err := decodePublic(publ)
	if err != nil {
		return nil, err
	}

	return &Tokenizer{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (t *Tokenizer) Generate(actor Actor) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":    actor.UserID.String(),
		"email":  actor.Email,
		"org_id": actor.OrgID.String(),
		"role":   actor.Role.String(),
		"iss":    TokenIssuer,
		"aud":    TokenAudience,
		"iat":    now.Unix(),
		"exp":    now.Add(TokenExpiration).Unix(),
	})

	signed, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (t *Tokenizer) Validate(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if _, ok := claims["sub"]; !ok {
		return nil, errors.New("missing subject")
	}

	if _, ok := claims["org_id"]; !ok {
		return nil, errors.New("missing Org ID")
	}

	if _, ok := claims["role"]; !ok {
		return nil, errors.New("missing role")
	}

	if _, ok := claims["iat"]; !ok {
		return nil, errors.New("missing issued at")
	}

	if _, ok := claims["exp"]; !ok {
		return nil, errors.New("missing expiration")
	}

	if _, ok := claims["iss"]; !ok {
		return nil, errors.New("missing issuer")
	}
	if claims["iss"] != TokenIssuer {
		return nil, errors.New("invalid issuer")
	}

	if _, ok := claims["aud"]; !ok {
		return nil, errors.New("missing audience")
	}
	if claims["aud"] != TokenAudience {
		return nil, errors.New("invalid audience")
	}

	return claims, nil
}

// decodePrivate decodes the private key from PEM format.
func decodePrivate(pemEncodedPriv string) (*ecdsa.PrivateKey, error) {
	blockPriv, _ := pem.Decode([]byte(pemEncodedPriv))
	x509EncodedPriv := blockPriv.Bytes
	return x509.ParseECPrivateKey(x509EncodedPriv)
}

// decodePublic decodes the public key from PEM format.
func decodePublic(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey, nil
}
