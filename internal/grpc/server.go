package grpc

import (
	"context"

	ptsv1 "github.com/JamshedJ/position-tracking-service/protos/gen/pts"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ptsv1.UnimplementedPositionTrackerServer
}

func Register(gRPC *grpc.Server) {
	ptsv1.RegisterPositionTrackerServer(gRPC, &serverAPI{})
}

func (s *serverAPI) GetNearbyPositions(*ptsv1.NearbyRequest, ptsv1.PositionTracker_GetNearbyPositionsServer) error {
	panic("implemet me")
}

func (s *serverAPI) UpdatePosition(context.Context, *ptsv1.PositionRequest) (*ptsv1.UpdateResponse, error) {
	panic("implemet me")
}