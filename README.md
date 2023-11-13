# Aggregator architecture

![Aggregator architecture](./docs/assets/architecture.drawio.svg)

Aggregator service provides a simple way to query information about a multi-cluster system. It launches a server that, when connected, returns available data about clusters and their states.

## Modules
### Server
Depending on the project, an HTTP or gRPC server is launched. It listens on the port AGGREGATOR_PORT passed as enviroment variable (default :8080) and calls the models module every time a client request is received. "AGGREGATOR_PORT" and "docker run --publish <port>" must be the same.

### Models
Data taxonomy is defined in *models.go* as Golang structs. 
The file *modeler.go* sends different queries (using *querier.go*) to Thanos and stores the received data. Once completed, it returns all the information in JSON format for web server and in the proto buffer defined format for gRPC server.

### Querier
Creates Prometheus API client and sends a query to Thanos. It retrieves the metrics and returns them in a response vector.


## Execution
Building the aggregator:
```bash
docker build . -t icos-aggregator
```

Launching the aggregator:
```bash
docker run -p 8080:8080 -e PROMETHEUS_ADDRESS=http://thanos.192.168.137.200.nip.io/ -e AGGREGATOR_PORT=8080 icos-aggregator
```

Connecting to the aggregator with HTTP server:
```bash
curl localhost:8080
```
  
Connecting to the aggregator with gRPC server (via Cognifog client test file):
```bash
export AGGREGATOR_PORT=8080 && go run test/protobuf/cognifog/server_client.go 
```
