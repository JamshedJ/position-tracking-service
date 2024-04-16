package main

import (
	"fmt"

	"github.com/JamshedJ/position-tracking-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger

	// TODO: init app

	// TODO: run GRPC server
}