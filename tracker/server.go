package tracker

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "P2P_BitTorrent/pb"
)

// Estructura para manejar la información del tracker.
type trackerServer struct {
	pb.UnimplementedTrackerServiceServer
	mu         sync.Mutex          // Para proteger el acceso concurrente a las estructuras.
	nodes      map[string]int      // Mapa de nodos activos con la cantidad de chunks que tienen.
	fileChunks map[string][]string // Mapa de archivos con la lista de nodos que tienen sus chunks.
}

// Crear una nueva instancia del servidor del tracker.
func NewTrackerServer() *trackerServer {
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
		return s.handlePut(fileName, req.FileSizeMb)
	}

	// Si la acción es 'get', gestionar la solicitud de descarga
	if action == "get" {
		return s.handleGet(fileName)
	}

	return &pb.JoinResponse{Message: "Acción desconocida."}, nil
}

// LeaveNetwork elimina un nodo de la red y actualiza la distribución de chunks.
func (s *trackerServer) LeaveNetwork(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeID := req.NodeId

	// Eliminar el nodo de fileChunks
	for chunkID, nodes := range s.fileChunks {
		// Filtrar los nodos que no sean el que está haciendo leave
		newNodes := []string{}
		for _, node := range nodes {
			if node != nodeID {
				newNodes = append(newNodes, node)
			}
		}

		// Si no quedan nodos almacenando el chunk, podemos eliminar el chunk (opcional)
		if len(newNodes) > 0 {
			s.fileChunks[chunkID] = newNodes
		} else {
			delete(s.fileChunks, chunkID)
		}

	}

	// Eliminar el nodo de la lista de nodos activos
	delete(s.nodes, nodeID)

	log.Printf("Nodo %s salió de la red y fue eliminado de todos los chunks.", nodeID)
	return &pb.LeaveResponse{Message: fmt.Sprintf("Nodo %s desconectado.", nodeID)}, nil
}
