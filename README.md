# ROS Project Backend

This is the backend service for the ROS Project frontend. It provides a RESTful API for managing ROS nodes, connections, and messages.

## Features

- Node management (create, retrieve)
- Connection management between nodes
- Message passing between nodes
- RESTful API endpoints
- CORS support for frontend integration
- JSON logging
- Graceful shutdown

## Prerequisites

- Go 1.21 or later
- Git

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/rosadsky/ros-project-backend.git
cd ros-project-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### Nodes
- `POST /api/ros/nodes` - Create a new node
- `GET /api/ros/nodes/:id` - Get a node by ID

### Connections
- `POST /api/ros/connections` - Create a connection between nodes

### Messages
- `POST /api/ros/messages` - Send a message
- `GET /api/ros/messages/:topic` - Get messages for a topic

## Configuration

The service can be configured using the `config/config.yaml` file. Available options:

- Server port and host
- Logging level and format
- CORS settings

## Development

To run the server in development mode with hot reload:

```bash
go run cmd/api/main.go
```

## Testing

To run the tests:

```bash
go test ./...
```

## License
