package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/EricNeid/go-webserver/database"
	"github.com/EricNeid/go-webserver/server"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	listenAddr string = ":5000"
	dbHost     string = "localhost"
	dbPort     int    = 5432
	dbUser     string = "postgres"
	dbPass     string = "postgres"
	dbName     string = "localdb"
)

func readEnvironmentVariables() {
	value, isSet := os.LookupEnv("LISTEN_ADDR")
	if isSet {
		listenAddr = value
	}

	value, isSet = os.LookupEnv("DB_HOST")
	if isSet {
		dbHost = value
	}

	value, isSet = os.LookupEnv("DB_PORT")
	if isSet {
		dbPort, _ = strconv.Atoi(value)
	}

	value, isSet = os.LookupEnv("DB_USER")
	if isSet {
		dbUser = value
	}

	value, isSet = os.LookupEnv("DB_PASS")
	if isSet {
		dbPass = value
	}

	value, isSet = os.LookupEnv("DB_NAME")
	if isSet {
		dbName = value
	}
}

func readCli() {
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.StringVar(&dbHost, "db-host", dbHost, "database host adress")
	flag.IntVar(&dbPort, "db-port", dbPort, "database port")
	flag.StringVar(&dbUser, "db-user", dbUser, "database user credential")
	flag.StringVar(&dbPass, "db-pass", dbPass, "database user password")
	flag.StringVar(&dbName, "db-name", dbName, "database name")

	flag.Parse()
}

func main() {
	readEnvironmentVariables()
	readCli()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	// create db pool
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("Connecting to db using: %s", dbUrl)
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		logger.Fatalf("Could not create database pool: %v\n", err)
	}
	createTables(logger, db)

	// create server
	server := server.NewApplicationServer(logger, db, listenAddr)
	go server.GracefullShutdown(quit, done)

	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	db.Close()
	logger.Println("Server stopped")
}

func createTables(logger *log.Logger, db *pgxpool.Pool) {
	err := database.CreateTablePositions(logger, db)
	if err != nil {
		logger.Fatalf("Could not create table %v\n", err)
	}

	err = database.CreateTableUsers(logger, db)
	if err != nil {
		logger.Fatalf("Could not create table %v\n", err)
	}
}
