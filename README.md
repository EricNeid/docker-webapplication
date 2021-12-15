# About

A simple REST-Full service. A dockerize environment with PostGIS database is provided as well.

## Components

* webserver
* PostGIS Database

## Quickstart

Use the provided Makefile for your convenience:

```bash
make all
```

Send test request:

```bash
curl -d '{"username":"max"}' -H "Content-Type: application/json" -X POST http://localhost:5000/users
```

```bash
curl -d '{"timestamp":"2021-06-15T09:00:00Z", "position": { "type": "Point", "coordinates": [20,30]}}' -H "Content-Type: application/json" -X POST http://localhost:5000/vehicleStates
```

## Testing

Unit and integration test (using a PostGIS Container) are provided. Running integration tests requires docker in your path.

To run unit tests:

```bash
go test -short ./...
```

To run all tests:

```bash
go test ./...
```
