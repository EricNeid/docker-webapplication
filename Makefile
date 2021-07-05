all: clean build run

build:
	docker-compose build --no-cache --force

clean:
	docker-compose down --remove-orphans

run:
	docker-compose up -d db
	sleep 10
	docker-compose up -d webserver

test:
	go test -short ./...

check:
	go test ./...
