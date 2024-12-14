# Kafka Implementation in GO

In this small repo we are creating two services ServiceA and serviceB

serviceA - Producer\
serviceB - Consumer



## Project structure

- **project-root/**
  - **.gitignore**
  - **docker-compose.yml**
  - **serviceA/**  
    - **Dockerfile**
    - **main.go**
    - **go.mod**
    - **go.sum**
  - **serviceB/** 
    - **Dockerfile**
    - **main.go**
    - **go.mod**
    - **go.sum**
  - **README.md**
  - **LICENSE**
## Working

In this docker-compose contains the config how the services communicates with one another\
serviceA and serviceB starts after kafka have started an ready.\
serviceA produces message inside the topic and serviceB reads from it by starting a go routine which reads\
from it asynchronously.

## Docker Compose Configuration

This project uses Docker Compose to manage the services. Below is the configuration for the services used in this project:

```yaml
version: '3.8'

services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"  # Expose zookeeper on port 2181

  kafka:
    image: wurstmeister/kafka
    ports:
      - "9092:9092"  # Expose Kafka on port 9092 for outside access
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181  # Connect Kafka to Zookeeper
      KAFKA_ADVERTISED_LISTENERS: INSIDE://kafka:9093,OUTSIDE://localhost:9092  # Advertise internal and external listeners with different ports
      KAFKA_LISTENERS: INSIDE://0.0.0.0:9093,OUTSIDE://0.0.0.0:9092  # Kafka will listen on different ports for inside and outside access
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "9093"]  # Check if Kafka's internal port 9093 is open and accessible
      interval: 30s  # Check every 30 seconds
      retries: 5  # Retry 5 times before considering Kafka as unhealthy
      timeout: 10s  # Timeout for each check is 10 seconds
      start_period: 10s  # Kafka will have 10 seconds to start before health checks begin

  service_a:
    build:
      context: ./serviceA
      dockerfile: Dockerfile
    container_name: service_a
    depends_on:
      kafka:
        condition: service_healthy  # Wait for Kafka to be available

  service_b:
    build:
      context: ./serviceB
      dockerfile: Dockerfile
    container_name: service_b
    depends_on:
      kafka:
        condition: service_healthy  # Wait for Kafka to be available


## Run Locally

Clone the project

```bash
  git clone git@github.com:piyushsharma67/kafka_implementation.git
```

Go to the project directory

```bash
  cd kafka_implementation
```

Start the project

```bash
  docker compose up --build
```

