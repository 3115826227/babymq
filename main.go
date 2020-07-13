package main

import (
	"github.com/3115826227/babymq/config"
	"github.com/3115826227/babymq/internel/communication/grpc"
)

func main() {

	config.InitConfig()

	grpcProvider := grpc.GetGrpcProvider()
	grpcProvider.Start()
}
