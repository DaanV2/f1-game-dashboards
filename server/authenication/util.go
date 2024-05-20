package authenication

import (
	"errors"

	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/DaanV2/f1-game-dashboards/server/users"
)

// Authenticator is the authentication service
func (a *Authenticator) ExtractUser(token *jwt.Token) (*users.User, error) {
	// Extract the user from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	// Get the user
	var (
		err  error
		user users.User
	)

	if user.Id, ok = claims["sub"].(string); !ok {
		err = errors.Join(err, errors.New("missing user id"))
	}
	if user.Email, ok = claims["email"].(string); !ok {
		err = errors.Join(err, errors.New("missing email"))
	}
	if user.Admin, ok = claims["admin"].(bool); !ok {
		err = errors.Join(err, errors.New("missing admin"))
	}
	if user.Guest, ok = claims["guest"].(bool); !ok {
		err = errors.Join(err, errors.New("missing guest"))
	}

	return &user, err
}

// ExtractGrant extracts the grant from the token
func (a *Authenticator) ExtractGrant(token *jwt.Token) (string, error) {
	// Extract the grant from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}

	// Get the grant
	grant, ok := claims["grant"].(string)
	if !ok {
		return "", errors.New("missing grant")
	}

	return grant, nil
}