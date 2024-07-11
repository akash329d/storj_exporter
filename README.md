# Storj Exporter for Prometheus

A Go-based exporter to pull information from Storj node APIs and export it for Prometheus monitoring. Supports monitoring multiple nodes, with each metric labeled by `node_id`. The existing python based exporter was pretty heavy on memory and additionally is a bit out of date with the current API responses.

## Features

- **Multi-node Support:** Collect metrics from multiple Storj nodes.
- **Dynamic Node Configuration:** Specify each node via environment variables.
- **Labelled Metrics:** Metrics are exported with a `node_id` label.

## Docker Usage

### Docker Hub Repository

Pull the Docker image from [Docker Hub](https://hub.docker.com/r/akash329d/storj_exporter).

### Example: Running the Docker Container

```sh
docker run -d \
  --name=StorjExporter \
  -e STORJ_NODE_1_URL=http://192.168.1.10:14002/ \
  -p 8000:8000 \
  akash329d/storj_exporter
```

### Docker Parameters

| Parameter         | Description                                    | Default Value |
|-------------------|------------------------------------------------|---------------|
| `EXPORTER_PORT`   | Port for the metrics server.                   | 8000          |
| `STORJ_NODE_%d_URL` | URL of a Storj node, %d starts at 1.           | N/A           |

## Accessing Metrics

Access the metrics at:
```arduino
http://<your-docker-host>:8000/metrics
```

Replace <your-docker-host> with the IP address or hostname of the machine running the Docker container.
