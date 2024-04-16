package grpc

import (
	"context"
	"math"
	"sync"

	"github.com/JamshedJ/position-tracking-service/internal/models"
	ptsv1 "github.com/JamshedJ/position-tracking-service/protos/gen/pts"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ptsv1.UnimplementedPositionTrackerServer
	positions map[string]models.ClientPosition
	rw        sync.RWMutex
}

func Register(gRPC *grpc.Server) {
	ptsv1.RegisterPositionTrackerServer(gRPC, &serverAPI{positions: make(map[string]models.ClientPosition)})
}

func (s *serverAPI) GetNearbyPositions(req *ptsv1.NearbyRequest, position ptsv1.PositionTracker_GetNearbyPositionsServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	centerLatitude := req.GetLatitude()
	centerLongitude := req.GetLongitude()
	radius := req.GetRadius()

	for _, pos := range s.positions {
		distance := haversine(centerLatitude, centerLongitude, pos.Latitude, pos.Longitude)
		if distance <= radius {
			// Отправить позицию клиента в поток
			err := position.Send(&ptsv1.PositionResponse{
				ClientId:  pos.ID,
				Latitude:  pos.Latitude,
				Longitude: pos.Longitude,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *serverAPI) UpdatePosition(ctx context.Context, req *ptsv1.PositionRequest) (*ptsv1.UpdateResponse, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	clientID := req.GetClientId()
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()

	s.positions[clientID] = models.ClientPosition{
		ID: clientID,
		Latitude: latitude,
		Longitude: longitude,
	}

	return &ptsv1.UpdateResponse{}, nil
}

// haversine вычисляет расстояние между двумя точками на поверхности Земли с использованием формулы Хаверсина.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Радиус Земли в километрах

	// Преобразовать координаты в радианы
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Разница в координатах
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Вычисление расстояния с использованием формулы Хаверсина
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}
