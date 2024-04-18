package grpc

import (
	"context"
	"fmt"

	"github.com/JamshedJ/position-tracking-service/internal/service"
	ptsv1 "github.com/JamshedJ/position-tracking-service/protos/gen/pts"
)

type Server struct {
	Service service.Service
	ptsv1.UnimplementedPositionTrackerServer
}

func NewServer(service service.Service) *Server {
	return &Server{
		Service: service,
	}
}

func (s *Server) GetNearbyPositions(req *ptsv1.NearbyRequest, stream ptsv1.PositionTracker_GetNearbyPositionsServer) error {
	results, err := s.Service.GetNearbyPositions(stream.Context(), req)
	if err != nil {
		return fmt.Errorf("error fetching nearby positions: %v", err)

	}

	for _, result := range results {
		if err := stream.Send(result); err != nil {
			return fmt.Errorf("error sending position response: %v", err)
		}
	}

	return nil
}

func (s *Server) UpdatePosition(ctx context.Context, req *ptsv1.PositionRequest) (*ptsv1.UpdateResponse, error) {
	resp, err := s.Service.UpdatePosition(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error updating position: %v", err)
	}
	return resp, nil
}
