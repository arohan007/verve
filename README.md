This high-performance RESTful service is designed to process 10,000+ requests per second while maintaining accurate deduplication of incoming requests. It provides a simple and robust solution for tracking unique request counts over time and supports seamless integration with external systems for real-time data processing.

`Core Functionality`

`API Endpoint`:

`Path: /api/verve/accept`
Accepts two query parameters:
id (mandatory, integer): A unique identifier for the request.
endpoint (optional, string): A URL to send unique request counts.
Returns:
"ok" (HTTP 200) for successful processing.
"failed" (HTTP 499) in case of errors.

`Deduplication`:

Ensures uniqueness of requests based on the id parameter.
Deduplication is achieved through a distributed Redis cache. Cache is flushed after every minute.

`Logging and External Communication`:

Logs the count of unique id values received every minute.
If an endpoint is provided, sends the unique request count via a POST request to the specified URL.

`Key Features`

`Distributed Scalability`: Handles deduplication across multiple service instances behind a load balancer using Redis.

`Real-Time Streaming`: Sends unique request counts to a distributed streaming platform (Kafka) for downstream processing.

`Lightweight and Fast`: Built with the Gin framework in Golang to ensure low latency and high throughput.

This service is ideal for high-throughput environments requiring efficient deduplication, seamless logging, and integration with external systems.
