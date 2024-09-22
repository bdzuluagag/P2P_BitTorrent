package main

import (
	"log"
	"net"

	pb "P2P_BitTorrent/pb"
	"P2P_BitTorrent/tracker" // El paquete tracker contendrá la lógica del servidor

	"google.golang.org/grpc"
)

func main() {
	// Configurar el servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	s := grpc.NewServer()
	trackerServer := tracker.NewTrackerServer()
	pb.RegisterTrackerServiceServer(s, trackerServer)

	log.Println("Tracker corriendo en el puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al correr el servidor: %v", err)
	}
}
