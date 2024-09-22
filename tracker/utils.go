package tracker

import (
	"strings"
)

// selectNodesForChunk selecciona varios nodos basados en la disponibilidad (menos chunks).
func (s *trackerServer) selectNodesForChunk(numReplicas int) []string {
	var selectedNodes []string

	// Crear una lista temporal de nodos que no han sido seleccionados
	availableNodes := make(map[string]int)
	for node, count := range s.nodes {
		availableNodes[node] = count
	}

	// Seleccionar numReplicas nodos
	for len(selectedNodes) < numReplicas {

		var selectedNode string
		minChunks := int(^uint(0) >> 1) // Reiniciar el valor máximo en cada iteración

		// Buscar el nodo con menos chunks que no haya sido seleccionado aún
		for node, count := range availableNodes {
			if count < minChunks && !contains(selectedNodes, node) { // Evitar seleccionar el mismo nodo
				minChunks = count
				selectedNode = node
			}
		}

		// Si no se seleccionó ningún nodo, salir del ciclo para evitar el bucle infinito
		if selectedNode == "" {
			break
		}

		selectedNodes = append(selectedNodes, selectedNode)
		delete(availableNodes, selectedNode) // Quitar el nodo de la lista temporal
	}

	return selectedNodes
}

// contains verifica si un nodo ya está en la lista de nodos seleccionados
func contains(nodes []string, node string) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}


// isValidChunk valida que chunkID pertenezca exactamente al archivo solicitado
func isValidChunk(fileName, chunkID string) bool {
	// Verificar si chunkID comienza con fileName seguido de un guion y un número
	if !strings.HasPrefix(chunkID, fileName) {
		return false
	}

	// Asegurarse de que después del nombre del archivo haya un guion y un número
	suffix := strings.TrimPrefix(chunkID, fileName)
	return strings.HasPrefix(suffix, "-") && len(suffix) > 1 && isNumeric(suffix[1:])
}

// isNumeric verifica si una cadena es numérica
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}