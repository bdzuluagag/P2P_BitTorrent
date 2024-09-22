package node

import (
	"P2P_BitTorrent/pb"
	"fmt"
)

// findChunk busca un chunk específico en la lista de chunks por su ID
func FindChunk(chunks []pb.StoreChunkRequest, chunkID string) *pb.StoreChunkRequest {
	for _, chunk := range chunks {
		if chunk.ChunkId == chunkID {
			return &chunk
		}
	}
	return nil
}

// Convierte el tamaño del archivo de string a int
func ParseSize(size string) (int, error) {
	var fileSize int
	_, err := fmt.Sscanf(size, "%d", &fileSize)
	if err != nil {
		return 0, err
	}
	return fileSize, nil
}

// Función para crear chunks de datos simulados
func CreateChunks(fileName string, totalSize, chunkSize int) []pb.StoreChunkRequest {
	var chunks []pb.StoreChunkRequest
	numChunks := totalSize / chunkSize
	for i := 0; i < numChunks; i++ {
		chunkID := fmt.Sprintf("%s-%d", fileName, i+1)
		chunks = append(chunks, pb.StoreChunkRequest{
			ChunkId:   chunkID,
			ChunkData: []byte(fmt.Sprintf("Datos del chunk %s", chunkID)), // Datos simulados
		})
	}
	return chunks

}
