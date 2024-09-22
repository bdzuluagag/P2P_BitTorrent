package tracker

import (
	"fmt"
	"log"

	pb "P2P_BitTorrent/pb"
)

// handlePut fragmenta el archivo y distribuye los chunks entre varios nodos.
func (s *trackerServer) handlePut(fileName string, fileSizeMb int32) (*pb.JoinResponse, error) {
	chunks := int(fileSizeMb) // Suponiendo 1 chunk por MB
	chunkMap := make(map[string]*pb.ChunkInfo)
	replicas := 3 // Definimos que queremos 3 réplicas por chunk, por ejemplo

	// Distribuir los chunks según la disponibilidad de los nodos
	for i := 0; i < chunks; i++ {
		chunkID := fmt.Sprintf("%s-%d", fileName, i+1) // Ej. Shakira.mp3-1

		// Seleccionar nodos para replicar el chunk
		selectedNodes := s.selectNodesForChunk(replicas)

		// Asignar los nodos seleccionados al chunk
		for _, targetNode := range selectedNodes {
			s.fileChunks[chunkID] = append(s.fileChunks[chunkID], targetNode)
			s.nodes[targetNode]++
			log.Printf("Chunk %s asignado al nodo %s", chunkID, targetNode)
		}
		chunkMap[chunkID] = &pb.ChunkInfo{
			Nodes: selectedNodes, // Lista de nodos que almacenan este chunk
		}
	}

	return &pb.JoinResponse{Message: fmt.Sprintf("Archivo %s subido y fragmentado exitosamente.", fileName), ChunkMap: chunkMap}, nil
}

// handleGet responde con los nodos que tienen los chunks del archivo solicitado.
func (s *trackerServer) handleGet(fileName string) (*pb.JoinResponse, error) {
	chunkMap := make(map[string]*pb.ChunkInfo)

	// Recoger los chunks y sus nodos en una estructura temporal
	for chunkID, nodes := range s.fileChunks {
		if isValidChunk(fileName, chunkID) { // Filtrar los chunks del archivo específico
			chunkMap[chunkID] = &pb.ChunkInfo{
				Nodes: nodes, // Asignar la lista de nodos que almacenan este chunk
			}
		}
	}

	if len(chunkMap) == 0 {
		return &pb.JoinResponse{Message: "Archivo no encontrado en la red."}, nil
	}

	log.Printf("Chunks encontrados para el archivo %s: %v", fileName, chunkMap)
	return &pb.JoinResponse{
		Message:  fmt.Sprintf("Nodos encontrados para los chunks del archivo %s", fileName),
		ChunkMap: chunkMap, // Enviamos el mapa de chunks y nodos asociados
	}, nil
}
