# Base image for Kafka and Zookeeper
FROM confluentinc/cp-kafka:7.4.0

# Copy Kafka and Zookeeper configuration files
COPY kafka-config/server.properties /etc/kafka/server.properties
COPY kafka-config/zookeeper.properties /etc/kafka/zookeeper.properties

# Expose Zookeeper and Kafka ports
EXPOSE 2181 9092

# Start Kafka and Zookeeper
CMD ["/etc/confluent/docker/run"]
