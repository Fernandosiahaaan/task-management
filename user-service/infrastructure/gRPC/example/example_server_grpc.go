package main

import (
	"context"
	"log"
	server "user-service/infrastructure/gRPC"
)

func main() {
	// Parameter untuk server gRPC
	param := server.ParamServerGrpc{
		Ctx:  context.Background(),
		Port: "50051",
	}

	// Membuat koneksi server gRPC
	serverGrpc, err := server.NewConnect(param)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Mulai mendengarkan permintaan
	go serverGrpc.StartListen()
	select {} // Menunggu selamanya
}
