package authenication

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/DaanV2/f1-game-dashboards/server/users"
	"github.com/charmbracelet/log"
)

// Token authenticates a user and returns a jwt token
func (a *Authenticator) Token(ctx context.Context, header string) (string, error) {
	logger := log.FromContext(ctx).With("header", header)

	// If basic, then its email and password, and return a jwt token
	// If jwt, then its a jwt token, possibly with a refresh token
	// else its assumed to be a guest name
	if strings.HasPrefix(header, "Bearer ") {
		logger.Debug("Refreshing token")
		return a.refreshToken(ctx, header[7:])
	}
	if strings.HasPrefix(header, "Basic ") {
		logger.Debug("Authenticating user with basic token to jwt token")
		return a.basicToken(ctx, header[6:])
	}

	guest := users.User{
		Id:       "guest: " + header,
		Email:    "guest@guest.com",
		Password: "",
		Admin:    false,
		Guest:    true,
	}

	return a.jwtToken(&guest, "guest")
}

func (a *Authenticator) basicToken(ctx context.Context, basic string) (string, error) {
	logger := log.FromContext(ctx)
	logger.Info("Authenticating user with basic token to jwt token")

	// Decode the basic token
	decoded, err := base64.StdEncoding.DecodeString(basic)
	if err != nil {
		return "", err
	}

	// Split the email and password
	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		return "", errors.New("invalid basic token")
	}

	// Authenticate the user
	user, err := a.users.Authenticate(ctx, parts[0], parts[1])
	if err != nil {
		return "", err
	}

	return a.jwtToken(user, "password")
}

func (a *Authenticator) refreshToken(ctx context.Context, tokenStr string) (string, error) {
	logger := log.FromContext(ctx).With("token", tokenStr)
	logger.Info("Refreshing token")

	// Verify the token
	token, err := a.jwtManager.Verify(tokenStr)
	if err != nil {
		return "", err
	}

	user, err := a.ExtractUser(token)
	if err != nil {
		logger.Error("error refreshing token", "error", err)
		return "", err
	}

	// Skip user exists check for guest, but enforce for adminds
	if !user.Guest || user.Admin {
		// Check if the user exists
		compare, err := a.users.GetByEmail(user.Email)
		if err != nil || compare.Id != user.Id || compare.Email != user.Email {
			logger.Error("token the user was made for has been changed", "error", err)
			return "", errors.New("invalid user")
		}
		user = compare
	}

	return a.jwtToken(user, "refresh")
}

func (a *Authenticator) jwtToken(user *users.User, grant string) (string, error) {
	logger := log.With(
		"id", user.Id,
		"email", user.Email,
		"admin", user.Admin,
		"guest", user.Guest,
		"grant", grant,
	)

	logger.Info("Creating jwt token")
	claims := jwt.MapClaims{
		"sub":   user.Id,
		"email": user.Email,
		"admin": user.Admin,
		"guest": user.Guest,
		"grant": grant,
	}

	return a.jwtManager.Sign(claims)
}
