# Censys KV Store

A key-value store service built with Go for the Censys interview take-home exercise.

This service provides the following REST API endpoints to store, retrieve, and delete key-value pairs:

- `GET /kv/{key}` - Retrieve a key-value pair by key
- `POST /kv` - Create or update a key-value pair  
- `DELETE /kv/{key}` - Soft delete a key-value pair by key

The service uses in-memory storage with soft delete functionality (records are marked as deleted but not removed).

## Design Decisions

- The service soft deletes records (they are marked as deleted but not removed). I did it intentionally because I believe it is a more real-world scenario. It would be easy to change to a hard delete if needed.
- When overwriting a soft deleted record, the service will undelete it (by emptying the `deleted_at` field) and update the value. In a real-world scenario, this could be handled better to keep track of the history of the record.
- All error response bodies will return JSON objects with `error` as the key.

## Testing Instructions

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running this service directly (outside Docker)

1. Run the service:
```bash
# NOTE: The default port is 8080. You can change it by setting the PORT environment variable.
go mod tidy
go run main.go # Service will be available at http://localhost:8080 (or the port you set)
```

**Note:** You can also run the container with `docker-compose up -d`, but you will need to create a Docker network first.
```bash
docker network create censys-network
docker-compose up -d
```

2. Test endpoints:
```bash
# Create/update a key-value pair
curl -X POST http://localhost:8080/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "my_key", "value": "my value"}'

# Get a key-value pair
curl http://localhost:8080/kv/my_key

# Soft delete a key-value pair
curl -X DELETE http://localhost:8080/kv/my_key
```

### Testing via the Test Client

The accompanying `kv-test-client` service is designed to test the `kv-store` service. Since both services run on containers, they need a docker network to communicate with each other.

1. Create the Docker network:
```bash
docker network create censys-network  # Only needed once
```

2. Run this service in Docker:
```bash
docker-compose up -d
```

3. Clone and run the test client (in a different terminal):
```bash
# In the desired directory
git clone git@github.com:augustoapg/censys-kv-test-client.git
cd kv-test-client
docker-compose up -d  # Test client will be available at http://localhost:8081
```

4. Call test client endpoints:
```bash
# Verifies that deleting a key-value pair works
curl -X POST http://localhost:8081/test_deletion

# Verifies that overwriting a key-value pair works
curl -X POST http://localhost:8081/test_overwrite
```

## Tearing down

To stop the services, run:
```bash
docker-compose down
```

To remove the network, run:
```bash
docker network rm censys-network
```
