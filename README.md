# P2P File Sharing System with gRPC, Replication & Fault Tolerance 🚀

Welcome to the **P2P File Sharing System**, a fully decentralized and distributed peer-to-peer (P2P) network built using **Go**, **gRPC**, and **Protobuf**. This project simulates the core functionality of the BitTorrent protocol, including file sharing with chunk replication, fault tolerance, and a central tracker.

## 🎯 Project Features

- **Decentralized P2P Network**: Each node in the network acts as both a client and a server, enabling efficient file sharing.
- **Tracker**: A central service that manages the list of nodes, tracks files, and stores information about which nodes hold chunks of each file.
- **Chunk-Based File Distribution**: Files are divided into chunks for efficient distribution across multiple nodes.
- **Replication**: Each chunk is replicated across multiple nodes to ensure availability and fault tolerance.
- **Fault Tolerance**: If a node goes offline, the file can still be reconstructed using the replicated chunks from other nodes.
- **gRPC Communication**: Nodes communicate via **gRPC**, ensuring efficient and scalable communication between peers and the tracker.
- **Built-in Commands**: Each node allows you to perform operations like uploading, downloading, and leaving the network through simple commands.

## 🛠️ Technologies Used

- **Go (Golang)**: Main programming language for the P2P system.
- **gRPC**: Enables efficient communication between peers and the tracker.
- **Protocol Buffers (Protobuf)**: For serializing structured data.
- **Amazon EC2**: For deploying and running the system in a real-world environment.
  
## 📁 Project Structure

```bash
.
├── tracker/                     # Tracker server that manages the nodes and file chunks
│   ├── server.go                # Tracker service implementation
│   ├── handlers.go              # Request handlers for the tracker
│   └── utils.go                 # Utility functions for the tracker
├── node/                        # Peer-to-peer nodes (client & server combined)
│   ├── server.go                # Server-side implementation of the node
│   └── utils.go                 # Utility functions for the node
├── proto/
│   └── peer.proto               # Protobuf definitions for the gRPC services
└── README.md                    # This README file
```

## 📦 Setup and Installation

### Prerequisites

- **Go 1.16+** installed on your machine.
- **gRPC** and **Protocol Buffers** tools installed.

### 1. Clone the Repository

```bash
git clone https://github.com/bdzuluagag/P2P_BitTorrent.git
cd P2P_BitTorrent
```

### 2. Install Dependencies

Install **gRPC** and **Protocol Buffers** plugins for Go:

```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### 3. Compile Protobuf Files

In the root of the project, compile the `.proto` files:

```bash
protoc --go_out=. --go-grpc_out=. proto/peer.proto
```

### 4. Run the Tracker

The tracker is the central service that manages the nodes and tracks which chunks of files are stored in each node.

```bash
cd cmd
go run tracker/tracker.go
```

The tracker will start on port `50051`.

### 5. Start Peer Nodes

Each node in the network can act as both a client and a server. Start a node on a specific port:

```bash
cd cmd
go run node/node.go
```

When prompted, enter a port number for the node (e.g., `50001`, `50002`).

### 6. Upload and Download Files

Each node can perform the following actions:

- **Put (Upload a file)**:
   ```bash
   put example.txt 10
   ```
   This will upload `example.txt` (which has a size of 10 MB), split it into chunks, and distribute it across available nodes.

- **Get (Download a file)**:
   ```bash
   get example.txt
   ```
   This will download all chunks of `example.txt` from the nodes, reconstruct the file, and store it locally.

- **Leave the network**:
   ```bash
   leave
   ```

## 🚀 Features Overview

### 1. **Join/Leave Network**
- When a node joins the network (through the `put` or `get` commands), it registers itself with the tracker. The tracker assigns chunks of files to nodes and updates its internal list.
- If a node leaves the network (via the `leave` command), the tracker removes it from the node list and updates its internal list.

### 2. **File Distribution and Replication**
- Files are split into chunks, and each chunk is replicated across multiple nodes to ensure redundancy.
- The default replication factor is 3, ensuring that each chunk is stored in 3 different nodes for fault tolerance.

### 3. **Fault Tolerance**
- If a node goes offline, other nodes that hold replicated chunks can serve the data.
- The tracker ensures that all file chunks remain available even if some nodes leave the network.

### 4. **gRPC Communication**
- Nodes communicate with each other and with the tracker using **gRPC** for efficient and scalable communication.
- All communication, including file uploads, downloads, and chunk transfers, is handled through gRPC requests and responses.

## 🧪 Example Usage

1. Start the tracker:

```bash
cd cmd
go run tracker/tracker.go
```

2. Start 3 peer nodes:

```bash
cd cmd
go run node/node.go
# Ingresar puerto: 50001

cd cmd
go run node/node.go
# Ingresar puerto: 50002

cd cmd
go run node/node.go
# Ingresar puerto: 50003
```

3. Upload a file from one of the nodes:

```bash
put shakira.mp3 10
```

4. Download the file from another node:

```bash
get shakira.mp3
```

## 🔧 Troubleshooting

- **Node Not Connecting:** Verify that the ports (`50000-50010`) are open and that the tracker is running on port `50051`.

## 🌟 Contributing

Feel free to submit issues or pull requests if you find any bugs or want to improve the project. Contributions are welcome!

---

🎉 **Thank you for checking out our P2P File Sharing System!** 🎉

---