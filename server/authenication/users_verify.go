package authenication

import (
	"context"

	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/DaanV2/f1-game-dashboards/server/users"
	"github.com/charmbracelet/log"
)

// Authenticator is a service that authenticates users
func (a *Authenticator) Verify(ctx context.Context, token string) (*jwt.Token, *users.User, error) {
	logger := log.FromContext(ctx).With("token", token)
	logger.Debug("Verifying token")

	// Verify the token
	t, err := a.jwtManager.Verify(token)
	if err != nil {
		logger.Warn("invalid token", "error", err)
		return t, nil, err
	}

	// Extract the user
	user, err := a.ExtractUser(t)
	if err != nil {
		logger.Warn("failed to extract user from token", "error")
		return t, nil, err
	}

	return t, user, nil
}