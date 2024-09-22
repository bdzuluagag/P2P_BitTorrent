package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"P2P_BitTorrent/node"
	pb "P2P_BitTorrent/pb"

	"google.golang.org/grpc"
)

const (
	trackerAddress = "localhost:50051" // Dirección y puerto del tracker
)

// Función principal del nodo
func main() {
	// Pedir al usuario que ingrese la ip:puerto del nodo
	fmt.Print("Ingrese la ip:puerto del nodo (ejemplo: localhost:50001, localhost:50002, ...): ")
	var nodePort string
	fmt.Scanln(&nodePort)

	// Inicia el servidor gRPC del nodo para manejar solicitudes de otros nodos
	go node.StartNodeServer(nodePort)

	// Conectar al tracker
	conn, err := grpc.Dial(trackerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar con el tracker: %v", err)
	}
	defer conn.Close()

	client := pb.NewTrackerServiceClient(conn)

	// Scanner para entrada de comandos del usuario
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Bienvenido al nodo cliente. Ingrese un comando:")
	fmt.Println("1. put [filename] [size_mb] - Para subir un archivo")
	fmt.Println("2. get [filename] - Para descargar un archivo")
	fmt.Println("3. leave - Para salir de la red")

	for scanner.Scan() {
		input := scanner.Text()
		commands := strings.Fields(input)

		if len(commands) == 0 {
			continue
		}

		switch commands[0] {
		case "put":
			if len(commands) != 3 {
				fmt.Println("Uso incorrecto. Ejemplo: put example.txt 10")
				continue
			}
			fileName := commands[1]
			fileSizeMb := commands[2]
			handlePut(client, fileName, fileSizeMb, nodePort)

		case "get":
			if len(commands) != 2 {
				fmt.Println("Uso incorrecto. Ejemplo: get example.txt")
				continue
			}
			fileName := commands[1]
			habdleGet(client, fileName, nodePort)

		case "leave":
			handleLeave(client, nodePort)
			return

		default:
			fmt.Println("Comando no reconocido. Intente de nuevo.")
		}
		fmt.Println("Ingrese otro comando:")
	}
}

// requestChunkFromNode simula la solicitud de un chunk desde otro nodo
func requestChunkFromNode(nodeAddress, chunkName string) {
	// Crear una conexión con el nodo destino
	conn, err := grpc.Dial(nodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("Error al conectar con el nodo %s: %v", nodeAddress, err)
		return
	}
	defer conn.Close()

	// Crear un cliente gRPC para el nodo
	client := pb.NewNodeServiceClient(conn)

	// Crear la solicitud del chunk
	req := &pb.ChunkRequest{
		ChunkId: chunkName, // Solicita el chunk específico
	}

	// Enviar la solicitud y recibir la respuesta
	res, err := client.RequestChunk(context.Background(), req)
	if err != nil {
		log.Printf("Error al solicitar chunk %s de %s: %v", chunkName, nodeAddress, err)
		return
	}

	log.Printf("Chunk %s recibido desde %s: %s", chunkName, nodeAddress, res.Message)
}

// Función para enviar un chunk a un nodo específico
func SendChunkToNode(nodeAddress string, chunk pb.StoreChunkRequest) {
	conn, err := grpc.Dial(nodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("Error al conectar con el nodo %s: %v", nodeAddress, err)
		return
	}
	defer conn.Close()

	client := pb.NewNodeServiceClient(conn)
	res, err := client.StoreChunk(context.Background(), &chunk)
	if err != nil {
		log.Printf("Error al enviar chunk %s a %s: %v", chunk.ChunkId, nodeAddress, err)
		return
	}

	log.Printf("Respuesta al enviar chunk %s a %s: %s", chunk.ChunkId, nodeAddress, res.Message)
}

// handlePut envía una solicitud para subir un archivo al tracker
func handlePut(client pb.TrackerServiceClient, fileName string, fileSizeMb string, nodeID string) {
	size, err := node.ParseSize(fileSizeMb)
	if err != nil {
		log.Printf("Error al parsear el tamaño del archivo: %v", err)
		return
	}

	// Crear la solicitud para el tracker
	req := &pb.JoinRequest{
		NodeId:     nodeID,
		Action:     "put",
		FileName:   fileName,
		FileSizeMb: int32(size),
	}

	// Enviar la solicitud al tracker
	res, err := client.JoinNetwork(context.Background(), req)
	if err != nil {
		log.Printf("Error al subir archivo: %v", err)
		return
	}

	fmt.Println(res.Message)

	// Simulación de chunks para enviar
	chunkSize := 1 // Suponiendo 1 MB por chunk
	chunks := node.CreateChunks(fileName, size, chunkSize)

	// Enviar cada chunk a los nodos correspondientes en el ChunkMap
	for chunkID, chunkInfo := range res.ChunkMap {
		chunk := node.FindChunk(chunks, chunkID)

		if chunk != nil {
			// Iterar sobre todos los nodos que almacenan este chunk
			for _, targetNode := range chunkInfo.Nodes {
				// Enviar el chunk al nodo correspondiente
				go SendChunkToNode(targetNode, *chunk)
			}
		}
	}
}

// habdleGet envía una solicitud para descargar un archivo al tracker y se conecta a los nodos correctos
func habdleGet(client pb.TrackerServiceClient, fileName string, nodeID string) {
	req := &pb.JoinRequest{
		NodeId:   nodeID,
		Action:   "get",
		FileName: fileName,
	}

	res, err := client.JoinNetwork(context.Background(), req)
	if err != nil {
		log.Printf("Error al descargar archivo: %v", err)
		return
	}

	fmt.Println(res.Message)

	// Solicitar cada chunk a su nodo correspondiente
	for chunkID, chunkInfo := range res.ChunkMap {
		// Seleccionamos un nodo para cada chunk (podríamos elegir el primer nodo o usar alguna estrategia)
		if len(chunkInfo.Nodes) > 0 {
			targetNode := chunkInfo.Nodes[0]             // Aquí seleccionamos el primer nodo como ejemplo
			go requestChunkFromNode(targetNode, chunkID) // Llamar a cada nodo concurrentemente con su chunk correspondiente
		} else {
			log.Printf("No hay nodos disponibles para el chunk %s", chunkID)
		}
	}
}

// handleLeave envía una solicitud para salir de la red al tracker
func handleLeave(client pb.TrackerServiceClient, nodeID string) {
	req := &pb.LeaveRequest{
		NodeId: nodeID,
	}

	res, err := client.LeaveNetwork(context.Background(), req)
	if err != nil {
		log.Printf("Error al salir de la red: %v", err)
		return
	}

	fmt.Println(res.Message)
}
