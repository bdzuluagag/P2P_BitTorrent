package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "P2P_BitTorrent/pb" // Importa el paquete generado de tu archivo .proto

	"google.golang.org/grpc"
)

// Estructura para manejar la información del tracker.
type trackerServer struct {
	pb.UnimplementedTrackerServiceServer
	mu         sync.Mutex          // Para proteger el acceso concurrente a las estructuras.
	nodes      map[string]int      // Mapa de nodos activos con la cantidad de chunks que tienen.
	fileChunks map[string][]string // Mapa de archivos con la lista de nodos que tienen sus chunks.
}

// Crear una nueva instancia del servidor del tracker.
func newTrackerServer() *trackerServer {
	return &trackerServer{
		nodes:      make(map[string]int),
		fileChunks: make(map[string][]string),
	}
}

// JoinNetwork permite a los nodos unirse a la red para subir o descargar archivos.
func (s *trackerServer) JoinNetwork(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeID := req.NodeId
	action := req.Action
	fileName := req.FileName

	// Verificar si el nodo ya está registrado
	if _, exists := s.nodes[nodeID]; !exists {
		s.nodes[nodeID] = 0 // Registrar nodo con 0 chunks inicialmente
		log.Printf("Nodo %s conectado a la red para acción: %s", nodeID, action)
	}

	// Si la acción es 'put', gestionar la subida y fragmentación del archivo
	if action == "put" {
		return s.handlePut(fileName, req.FileSizeMb, nodeID)
	}

	// Si la acción es 'get', gestionar la solicitud de descarga
	if action == "get" {
		return s.handleGet(fileName, nodeID)
	}

	return &pb.JoinResponse{Message: "Acción desconocida."}, nil
}

// handlePut fragmenta el archivo y distribuye los chunks entre los nodos.
func (s *trackerServer) handlePut(fileName string, fileSizeMb int32, nodeID string) (*pb.JoinResponse, error) {
	chunks := int(fileSizeMb) // Suponiendo 1 chunk por MB
	chunkNodes := []string{}

	// Distribuir los chunks según la disponibilidad de los nodos
	for i := 0; i < chunks; i++ {
		targetNode := s.selectNodeForChunk()
		chunkID := fmt.Sprintf("%s-%d", fileName, i+1) // Ej. Shakira.mp3-1
		s.fileChunks[chunkID] = append(s.fileChunks[chunkID], targetNode)
		s.nodes[targetNode]++
		chunkNodes = append(chunkNodes, targetNode)
		log.Printf("Chunk %s asignado al nodo %s", chunkID, targetNode)
	}

	return &pb.JoinResponse{Message: fmt.Sprintf("Archivo %s subido y fragmentado exitosamente.", fileName)}, nil
}

// selectNodeForChunk selecciona un nodo basado en la disponibilidad (menos chunks).
func (s *trackerServer) selectNodeForChunk() string {
	var selectedNode string
	minChunks := int(^uint(0) >> 1) // Iniciar con el máximo valor posible

	for node, count := range s.nodes {
		if count < minChunks {
			minChunks = count
			selectedNode = node
		}
	}

	return selectedNode
}

// handleGet responde con los nodos que tienen los chunks del archivo solicitado.
func (s *trackerServer) handleGet(fileName, nodeID string) (*pb.JoinResponse, error) {
	chunkNodes := []string{}

	// Buscar los chunks del archivo y sus respectivos nodos
	for chunkID, nodes := range s.fileChunks {
		if chunkID[:len(fileName)] == fileName { // Filtrar los chunks del archivo específico
			chunkNodes = append(chunkNodes, nodes...)
		}
	}

	if len(chunkNodes) == 0 {
		return &pb.JoinResponse{Message: "Archivo no encontrado en la red."}, nil
	}

	log.Printf("Nodos encontrados para el archivo %s: %v", fileName, chunkNodes)
	return &pb.JoinResponse{Message: fmt.Sprintf("Nodos encontrados: %v", chunkNodes)}, nil
}

// LeaveNetwork elimina un nodo de la red y actualiza la distribución de chunks.
func (s *trackerServer) LeaveNetwork(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeID := req.NodeId
	delete(s.nodes, nodeID)
	// Aquí podríamos implementar la redistribución de los chunks del nodo que se va.

	log.Printf("Nodo %s salió de la red.", nodeID)
	return &pb.LeaveResponse{Message: fmt.Sprintf("Nodo %s desconectado.", nodeID)}, nil
}

func main() {
	// Configurar el servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

	s := grpc.NewServer()
	tracker := newTrackerServer()
	pb.RegisterTrackerServiceServer(s, tracker)

	log.Println("Tracker corriendo en el puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al correr el servidor: %v", err)
	}
}
