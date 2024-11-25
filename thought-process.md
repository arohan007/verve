`README: High-Level Design and Thought Process for the Service`

`Overview`
This project implements a high-performance RESTful service designed to handle 10,000+ requests per second using Golang with the Gin framework. The service provides the following core features and extensions:

`Core Service`:
A single endpoint: /api/verve/accept.
Accepts a mandatory integer id and an optional HTTP endpoint URL (endpoint) as query parameters.
Logs unique requests (based on the id parameter) per minute.
Optionally sends the count of unique requests to an external endpoint.

`Extensions`:
Support for deduplication in distributed environments using Redis.
Sends unique request counts to a distributed streaming system (Kafka).
Fully containerized using Docker, including Redis, Kafka, and the service itself.

`High-Level Design`
1. `Service Architecture`
`Core API`
The /api/verve/accept endpoint is built using the Gin framework for high performance.
Parameters:
id (integer): Mandatory, uniquely identifies the request.
endpoint (string): Optional, URL to which the service sends the count of unique requests.
Response:
200 OK: For successful requests with "ok".
400 Error: For failed requests with "failed".

`Unique ID Deduplication`
A deduplication mechanism tracks unique id values for each minute.
In a single-instance deployment:
A Go map is used to maintain state.
Data is flushed and logged every minute.
In a multi-instance (load-balanced) deployment:
Redis is used as a shared distributed data store for deduplication.

`Logging and External Communication`
Logs the count of unique id values per minute using Go's log package.
If an endpoint parameter is provided:
Sends a POST request with the unique count and logs the response status code.


2. `Extensions`
`Extension`: `Deduplication Across Load Balancers`
Redis is introduced as a distributed in-memory store.
Each instance of the service writes the unique id into a Redis key with a 1-minute expiry (SET with TTL).
Deduplication is achieved using Redis' atomic operations (SADD for sets).
`Extension`: `Kafka Integration`
Kafka is used to decouple processing from logging and real-time monitoring.
Each service instance produces messages to a Kafka topic (unique-request-count).
`Extension`: `Fully Containerized Deployment`
Docker is used for creating containers for the service and dependencies (Redis, Kafka).
Docker Compose simplifies managing the multi-container setup.

`Design Choices`
`Framework: Gin`
Chosen for its lightweight and high-performance HTTP handling.
Ideal for achieving the 10k requests/second benchmark.

`Concurrency: Go Routines`
Go routines allow handling multiple requests simultaneously, minimizing latency.

`Data Store: Redis`
Redis is used for distributed deduplication because of its:
In-memory data store for high-speed operations.
Atomic operations for accurate deduplication.
Built-in TTL support for efficient memory management.

`Streaming: Kafka`
Kafka decouples the service's processing and external consumers, ensuring scalability.
Reliable, high-throughput distributed system for real-time data processing.

`Containerization: Docker`
Simplifies deployment and dependency management.
Ensures consistent environment across development, testing, and production.

`Deployment`
`Using Docker Compose`
Navigate to docker directory. (For eg: `/home/arohan/go/src/projects/verve/docker`)
`Build and Start the Service:`
```
docker-compose up --build
```
`Access the Service`:
API available at http://localhost:8080/api/verve/accept.
