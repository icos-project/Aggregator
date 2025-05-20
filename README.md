# Aggregator architecture

![Aggregator architecture](./docs/assets/architecture.drawio.svg)

Aggregator service provides a simple way to query information about a multi-cluster system. It launches a server that, when connected, returns available data about clusters and their states.

## Interfaces documentation

### HTTP
Topology information and infrastructure data is served via HTTP following the data structure defined in the Swagger documentation located at [./docs](./docs).

## Modules
### Server
The HTTP server listens on the provided port passed as environment variable and calls the models module every time a client request is received.

### Models
Data taxonomy is defined in *models.go* as Golang structs. 
The file *modeler.go* sends different queries (using *querier.go*) to Thanos and stores the received data. Once completed, it returns all the information in JSON format.

### Querier
Creates Prometheus API client and sends a query to Thanos. It retrieves the metrics and returns them in a response vector.


## Execution

#### Building the aggregator:
```bash
docker build . -t icos-aggregator
```
  
#### Launching the aggregator:

Environment variables:
- PROMETHEUS_ADDRESS: The address where Prometheus/Thanos is located.  
- HTTP_PORT: Port where the HTTP server will listen. Defaults to 8080 if not set.
- KEY: If set, Keycloak is enabled and public key set as KEY.

```bash
docker run -p 8080:8080 -e PROMETHEUS_ADDRESS=http://thanos.192.168.137.200.nip.io/ -e HTTP_PORT=8080 icos-aggregator
```
(Note: HTTP_PORT must be published with the option -p to be able to run the container correctly)

Current Key:
```bash
KEY=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgTGF4mKVEa+eWX0S/+EWIfkkqbLba5WuQ1KKGRQz+P56Y0WNRbgjNl0CObndffmixbpgp4kg5jKq78HoFFP7bj0jQSNC3P26K9xPolFXbAlNJe41VMdI7xOkOF0D9GCplEylGlUlCgpaBnbloI4WcbH+RQ6n6Qp6MmNE+/xC3OMMhgEBacbiGtIR71N/HcDYDUORE335sSRpkrHhMxk3eWgZdIyfX88n9UkI3CtgNGIGgF8/w7ZYF2XBmVuv5+QE9d5fM9pZKWQnzBnsMJy4Xc+qZrZMI45KCHIW/DSFVGSsGboiVHSNVOu3mNhPSjvJtIH/7lItCG6m5zvBAvNf8QIDAQAB
```

#### Connecting to the aggregator:
```bash
curl localhost:8080
```

# Legal
The Aggregator is released under the Apache 2.0 license.
Copyright © 2022-2024 Bull SAS. All rights reserved.

🇪🇺 This work has received funding from the European Union's HORIZON research and innovation programme under grant agreement No. 101070177.
