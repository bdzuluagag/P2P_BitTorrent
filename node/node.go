package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "P2P_BitTorrent/pb" // Asegúrate de que esta ruta es correcta para tus archivos generados

	"google.golang.org/grpc"
)

const (
	trackerAddress = "localhost:50051" // Dirección y puerto del tracker
)

// Función principal del nodo
// Función principal del nodo
func main() {
	// Pedir al usuario que ingrese el ID del nodo
	fmt.Print("Ingrese el ID del nodo (ejemplo: node1, node2, ...): ")
	var nodeID string
	fmt.Scanln(&nodeID)

	// Establecer conexión con el tracker
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
			handlePut(client, fileName, fileSizeMb, nodeID)

		case "get":
			if len(commands) != 2 {
				fmt.Println("Uso incorrecto. Ejemplo: get example.txt")
				continue
			}
			fileName := commands[1]
			handleGet(client, fileName, nodeID)

		case "leave":
			handleLeave(client, nodeID)
			return

		default:
			fmt.Println("Comando no reconocido. Intente de nuevo.")
		}
		fmt.Println("Ingrese otro comando:")
	}
}

// Actualizamos handlePut y handleGet para incluir el nodoID
func handlePut(client pb.TrackerServiceClient, fileName string, fileSizeMb string, nodeID string) {
	// Convertir tamaño del archivo a entero
	size, err := parseSize(fileSizeMb)
	if err != nil {
		log.Printf("Error al parsear el tamaño del archivo: %v", err)
		return
	}

	// Crear la solicitud para el tracker
	req := &pb.JoinRequest{
		NodeId:     nodeID, // Usar el ID específico del nodo
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
}

func handleGet(client pb.TrackerServiceClient, fileName string, nodeID string) {
	req := &pb.JoinRequest{
		NodeId:   nodeID, // Usar el ID específico del nodo
		Action:   "get",
		FileName: fileName,
	}

	// Enviar la solicitud al tracker
	res, err := client.JoinNetwork(context.Background(), req)
	if err != nil {
		log.Printf("Error al descargar archivo: %v", err)
		return
	}

	fmt.Println(res.Message)
}

func handleLeave(client pb.TrackerServiceClient, nodeID string) {
	req := &pb.LeaveRequest{
		NodeId: nodeID,
	}

	// Enviar la solicitud de salida al tracker
	res, err := client.LeaveNetwork(context.Background(), req)
	if err != nil {
		log.Printf("Error al salir de la red: %v", err)
		return
	}

	fmt.Println(res.Message)
}

// Convierte el tamaño del archivo de string a int
func parseSize(size string) (int, error) {
	var fileSize int
	_, err := fmt.Sscanf(size, "%d", &fileSize)
	if err != nil {
		return 0, err
	}
	return fileSize, nil
}
