all: clean build start

build:
	docker-compose build --no-cache --force

clean:
	docker-compose down --remove-orphans

start:
	docker-compose up -d db
	sleep 10
	docker-compose up -d webserver

test:
	go test -short ./...

test_full:
	go test ./...

check:
	go test ./...
