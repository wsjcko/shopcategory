package main

import (
	pb "github.com/wsjcko/shopcategory/proto/pb"
	"shopcategory/handler"

	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
)

var (
	service = "shopcategory"
	version = "latest"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
	)
	srv.Init()

	// Register handler
	pb.RegisterShopcategoryHandler(srv.Server(), new(handler.Shopcategory))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}