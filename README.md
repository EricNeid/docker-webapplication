# About

A simple REST-Full service. A dockerized environment with PostGIS database is provided aswell.

## Components

* webserver
* PostGIS Database

## Quickstart

Use the provided Makefile for your convience:

```bash
make all
```

Send test request:

```bash
curl -d '{"username":"max"}' -H "Content-Type: application/json" -X POST http://localhost:5000/user
```

## Testing

Unit and intengration test (using a PostGIS Container) are provided.

To run unit tests:

```bash
go test -short ./...
```

To run all tests:

```bash
go test ./...
```
