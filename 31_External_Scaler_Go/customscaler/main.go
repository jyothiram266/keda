package main

import (
	"log"
	"net"
	pb "customscaler/externalscaler"

	"google.golang.org/grpc"
)

func main() {
	go RunManagementApi()
	go reduceCustomQueueLength()

	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")
	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})

	log.Println("Listening external scaler on :6000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
