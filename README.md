# Micro-Batching Service

## Build and Run

```bash
make codegen  # generate code from openapi.yaml
go build .
go run .
```

## Usage

Configurations can be set in `config.json` file, can also be updated via endpoints when the program runs

### Set Frequency

To call BatchProcessor every 5 seconds via the `/batch-frequency` endpoint:

```bash
curl http://localhost:8080/batch-frequency -H "Content-Type:application/json" -d '{"frequency":10}'
```

### Get Frequency

To get the current frequency of BatchProcessor:

```bash
curl http://localhost:8080/batch-frequency
```

### Set Batch Size

To set the batch size of BatchProcessor via the `/batch-size` endpoint:

```bash
curl http://localhost:8080/batch-size -H "Content-Type:application/json" -d '{"batch-size":10}'
```

### Get Batch Size

To get the current batch size of BatchProcessor:

```bash
curl http://localhost:8080/batch-size
```

### Add Job to Batching Service


## Design

### Preprocessing

This service can do an optional preprocessing step before sending the data to BatchProcessor. The preprocessing is disabled by default, can be turned on via `preprocess` endpoint. The preprocessing takes a list of jobs, and returns processed jobs. The idea is to allow possible filtering or merging to reduce the number of jobs to be processed by BatchProcessor.

For example, 
