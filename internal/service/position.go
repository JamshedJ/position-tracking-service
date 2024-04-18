package service

import (
	"context"
	"math"

	"github.com/JamshedJ/position-tracking-service/internal/dto"
	mdb "github.com/JamshedJ/position-tracking-service/internal/storage/mongodb"
	"github.com/JamshedJ/position-tracking-service/protos/gen/pts"
)

type Service interface {
	GetNearbyPositions(ctx context.Context, req *pts.NearbyRequest) ([]*pts.PositionResponse, error)
	UpdatePosition(ctx context.Context, req *pts.PositionRequest) (*pts.UpdateResponse, error)
}

func NewService(storage mdb.Storage) Service {
	return &service{
		storage: storage,
	}
}

type service struct {
	storage mdb.Storage
}

func (s *service) GetNearbyPositions(ctx context.Context, req *pts.NearbyRequest) ([]*pts.PositionResponse, error) {
	degreeDistance := req.Radius / 111.0
	squareDegreeDistance := degreeDistance * degreeDistance
	minLatitude := req.Latitude - degreeDistance
	maxLatitude := req.Latitude + degreeDistance
	minLongitude := req.Longitude - degreeDistance/math.Cos(req.Latitude*math.Pi/180)
	maxLongitude := req.Longitude + degreeDistance/math.Cos(req.Latitude*math.Pi/180)

	posRes, err := s.storage.GetNearbyPositions(ctx, req.Latitude, req.Longitude, dto.PositionParams{
		SquareDegreeDistance: squareDegreeDistance,
		MinLatitude:          minLatitude,
		MaxLatitude:          maxLatitude,
		MinLongitude:         minLongitude,
		MaxLongitude:         maxLongitude,
	})
	if err != nil {
		return nil, err
	}

	return posRes, nil
}

func (s *service) UpdatePosition(ctx context.Context, req *pts.PositionRequest) (*pts.UpdateResponse, error) {
	err := s.storage.UpdatePosition(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pts.UpdateResponse{
		Success: true,
		Message: "Position updated successfully",
	}, nil
}
