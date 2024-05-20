package authenication

import (
	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/DaanV2/f1-game-dashboards/server/users"
)

type Authenticator struct {
	users      *users.UserManagement
	jwtManager *jwt.JwtService
}

func NewAuthenticator(users *users.UserManagement, jwtManager *jwt.JwtService) *Authenticator {
	return &Authenticator{
		users:      users,
		jwtManager: jwtManager,
	}
}