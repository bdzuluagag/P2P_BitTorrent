package node

import (
	"P2P_BitTorrent/pb"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// Estructura del nodo para manejar tanto el servidor como el cliente gRPC
type nodeServer struct {
	pb.UnimplementedNodeServiceServer
	mu     sync.Mutex
	chunks map[string][]byte // Mapa para almacenar los chunks del nodo
}

// Inicializar el servidor con un mapa de chunks vacío
func newNodeServer() *nodeServer {
	return &nodeServer{
		chunks: make(map[string][]byte),
	}
}

// startNodeServer inicia el servidor gRPC del nodo
func StartNodeServer(nodeID string) {

	// Separar la IP del puerto
	_, port, err := net.SplitHostPort(nodeID)
	if err != nil {
		log.Fatalf("Error al separar la dirección IP y el puerto: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port)) // Asigna un puerto específico para cada nodo basado en su ID
	if err != nil {
		log.Fatalf("Error al iniciar el servidor del nodo: %v", err)
	}

	s := grpc.NewServer()
	node := newNodeServer()
	pb.RegisterNodeServiceServer(s, node)

	log.Printf("Nodo escuchando en %s...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al correr el servidor del nodo: %v", err)
	}
}

// RequestChunk maneja la solicitud de un chunk desde otro nodo
func (s *nodeServer) RequestChunk(ctx context.Context, req *pb.ChunkRequest) (*pb.ChunkResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chunkID := req.ChunkId
	data, exists := s.chunks[chunkID]

	if !exists {
		log.Printf("El chunk %s no está disponible en este nodo", chunkID)
		return &pb.ChunkResponse{
			Message: fmt.Sprintf("El chunk %s no está disponible", chunkID),
		}, nil
	}

	log.Printf("Solicitud recibida para el chunk %s", chunkID)
	return &pb.ChunkResponse{
		ChunkData: data,
		Message:   fmt.Sprintf("Chunk %s enviado correctamente", chunkID),
	}, nil
}

// Función para manejar la solicitud de almacenamiento de un chunk
func (s *nodeServer) StoreChunk(ctx context.Context, req *pb.StoreChunkRequest) (*pb.StoreChunkResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.chunks[req.ChunkId] = req.ChunkData
	log.Printf("Chunk %s almacenado correctamente en el nodo", req.ChunkId)
	return &pb.StoreChunkResponse{
		Message: fmt.Sprintf("Chunk %s almacenado correctamente", req.ChunkId),
	}, nil

}
