package api

import (
	"context"

	grpc_gen "github.com/DaanV2/f1-game-dashboards/server/api/grpc"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/charmbracelet/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ grpc_gen.ChairServiceServer = &grpcServer{}

// CreateChair implements grpc_gen.ChairServiceServer.
func (s *grpcServer) CreateChair(ctx context.Context, req *grpc_gen.CreateChairRequest) (*grpc_gen.CreateChairResponse, error) {
	response := grpc_gen.CreateChairResponse{}
	logger := log.FromContext(ctx)
	if _, err := s.mustBeAdmin(ctx); err != nil {
		return nil, err
	}

	c := req.GetChair()
	if c == nil {
		return &response, status.Error(codes.InvalidArgument, "chair is required")
	}
	logger = logger.With(
		"name", c.GetName(),
		"port", c.GetPort(),
		"active", c.GetActive(),
	)
	requestChair := chairFromProto(c)

	logger.Info("checking if chair exists")
	if _, exists := s.chairs.Get(requestChair.Id()); exists {
		logger.Warn("chair already exists")
		return &response, status.Error(codes.AlreadyExists, "chair already exists")
	}

	logger.Info("adding chair")
	s.chairs.Add(requestChair)
	response.Chair = chairToProto(requestChair)
	return &response, nil
}

// DeleteChair implements grpc_gen.ChairServiceServer.
func (s *grpcServer) DeleteChair(ctx context.Context, req *grpc_gen.DeleteChairRequest) (*grpc_gen.DeleteChairResponse, error) {
	response := grpc_gen.DeleteChairResponse{}
	port := req.GetPort()
	logger := log.FromContext(ctx).With("port", port)
	if _, err := s.mustBeAdmin(ctx); err != nil {
		return nil, err
	}

	logger.Info("getting chair")
	if port == "" || !sessions.IsChairId(port) {
		logger.Error("port is required")
		return &response, status.Error(codes.InvalidArgument, "port is required")
	}

	log.Info("deleting chair", "port", port)
	if _, exists := s.chairs.Get(port); !exists {
		logger.Info("chair not found")
		return &response, status.Error(codes.NotFound, "chair not found")
	}

	s.chairs.Remove(port)

	return &response, nil
}

// GetChair implements grpc_gen.ChairServiceServer.
func (s *grpcServer) GetChair(ctx context.Context, req *grpc_gen.GetChairRequest) (*grpc_gen.GetChairResponse, error) {
	response := grpc_gen.GetChairResponse{}
	port := req.GetPort()
	logger := log.FromContext(ctx).With("port", port)
	if _, err := s.atleastGuest(ctx); err != nil {
		return nil, err
	}

	if port == "" || !sessions.IsChairId(port) {
		return &response, status.Error(codes.InvalidArgument, "port is required")
	}

	logger.Info("getting chair")
	chair, exists := s.chairs.Get(port)
	if !exists {
		logger.Info("chair not found")
		return &response, status.Error(codes.NotFound, "chair not found")
	}

	response.Chair = chairToProto(chair)
	return &response, nil
}

// UpdateChair implements grpc_gen.ChairServiceServer.
func (s *grpcServer) UpdateChair(ctx context.Context, req *grpc_gen.UpdateChairRequest) (*grpc_gen.UpdateChairResponse, error) {
	logger := log.FromContext(ctx)
	response := grpc_gen.UpdateChairResponse{}
	if _, err := s.atleastGuest(ctx); err != nil {
		return nil, err
	}

	c := req.GetChair()
	if c == nil {
		return &response, status.Error(codes.InvalidArgument, "chair is required")
	}
	requestChair := chairFromProto(c)
	logger = logger.With("port", requestChair.Port)

	oldChair, exists := s.chairs.Get(requestChair.Id())
	if !exists {
		logger.Info("chair not found")
		return &response, status.Error(codes.NotFound, "chair not found")
	}
	if oldChair.Port != requestChair.Port {
		logger.Info("ports do not match")
		return &response, status.Error(codes.InvalidArgument, "ports do not match")
	}

	updateChair := sessions.NewChair(
		requestChair.Name,
		oldChair.Port,
		requestChair.Active,
	)
	response.Chair = chairToProto(updateChair)

	// Must be admin to change name
	if oldChair.Name != updateChair.Name {
		if _, err := s.mustBeAdmin(ctx); err != nil {
			return nil, err
		}
	}

	s.chairs.Update(updateChair)

	return &response, nil
}

// ListChairs implements grpc_gen.ChairServiceServer.
func (s *grpcServer) ListChairs(ctx context.Context, req *grpc_gen.ListChairsRequest) (*grpc_gen.ListChairsResponse, error) {
	if _, err := s.atleastGuest(ctx); err != nil {
		return nil, err
	}
	chairs := s.chairs.All()

	chrs := make([]*grpc_gen.Chair, 0, len(chairs))
	response := grpc_gen.ListChairsResponse{
		Chairs: chrs,
	}

	for _, chair := range chairs {
		chrs = append(chrs, chairToProto(chair))
	}

	return &response, nil
}

func chairToProto(chair sessions.Chair) *grpc_gen.Chair {
	return &grpc_gen.Chair{
		Name:   chair.Name,
		Active: chair.Active,
		Port:   int32(chair.Port),
	}
}

func chairFromProto(chair *grpc_gen.Chair) sessions.Chair {
	return sessions.NewChair(chair.GetName(), int(chair.GetPort()), chair.GetActive())
}
