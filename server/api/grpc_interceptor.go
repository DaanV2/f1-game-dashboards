package api

import (
	"context"
	"errors"

	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/DaanV2/f1-game-dashboards/server/users"
	"github.com/charmbracelet/log"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	AuthenicationKey   struct{}
	AuthenicationValue struct {
		Token *jwt.Token
		User  *users.User
		Error error
	}
)

func (s *grpcServer) interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	logger := log.FromContext(ctx).With("method", info.FullMethod)
	authV := AuthenicationValue{
		Token: nil,
		User:  nil,
		Error: nil,
	}

	// Get the metadata from the context, possibly extract the JWT
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		auth := md.Get("authorization")
		if len(auth) > 0 {
			t, u, err := s.authenicator.Verify(ctx, auth[0])
			authV = AuthenicationValue{
				Token: t,
				User:  u,
				Error: err,
			}
			if u != nil {
				logger = logger.With(
					"user", u.Email,
					"admin", u.Admin,
					"guest", u.Guest,
				)
			}
			if t != nil {
				logger = logger.With("valid", t.Valid)
			}
			ctx = log.WithContext(ctx, logger)
		}
	}

	ctx = context.WithValue(ctx, AuthenicationKey{}, authV)

	// Next
	return handler(ctx, req)
}

func (s *grpcServer) getAuth(ctx context.Context) AuthenicationValue {
	v := ctx.Value(AuthenicationKey{})
	if v == nil {
		return AuthenicationValue{
			Error: errors.New("no authenication value found"),
		}
	}
	return v.(AuthenicationValue)
}

// mustBeAdmin returns the user if the user is an admin, otherwise an error is returned.
func (s *grpcServer) mustBeAdmin(ctx context.Context) (*users.User, error) {
	auth := s.getAuth(ctx)
	if auth.Error != nil {
		return auth.User, auth.Error
	}
	if auth.User.Guest || !auth.User.Admin {
		return auth.User, status.Error(codes.PermissionDenied, "user is not an admin")
	}
	return auth.User, nil
}

// atleastGuest returns the user if the user is authenicated atleast as a guest, otherwise an error is returned.
func (s *grpcServer) atleastGuest(ctx context.Context) (*users.User, error) {
	auth := s.getAuth(ctx)
	if auth.Error != nil {
		return auth.User, auth.Error
	}
	return auth.User, nil
}