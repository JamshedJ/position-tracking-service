package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/JamshedJ/position-tracking-service/internal/config"
	grpcSever "github.com/JamshedJ/position-tracking-service/internal/grpc"
	"github.com/JamshedJ/position-tracking-service/internal/service"

	mdb "github.com/JamshedJ/position-tracking-service/internal/storage/mongodb"
	ptsv1 "github.com/JamshedJ/position-tracking-service/protos/gen/pts"
)

func main() {
	cfg := config.MustLoad()

	collection := mdb.New(cfg.MongoDB.Uri, cfg.MongoDB.DBName, cfg.MongoDB.CollectionName)

	storage := mdb.NewStorage(collection)
	svc := service.NewService(storage)
	server := grpc.NewServer()
	ptsv1.RegisterPositionTrackerServer(server, &grpcSever.Server{Service: svc})

	listener, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	reflection.Register(server)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")
	server.GracefulStop()
}
