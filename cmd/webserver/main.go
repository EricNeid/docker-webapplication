package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/EricNeid/go-webserver/server"
	"github.com/gin-gonic/gin"
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

	logFile, _ := os.Create("server.log")
	logWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(logWriter, "main", log.LstdFlags)

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

	// create server
	gin.SetMode(gin.ReleaseMode)
	server := server.NewApplicationServer(db, listenAddr)
	server.SetLogWriter(logWriter)
	go server.GracefullShutdown(quit, done)

	logger.Println("Creating database structure", listenAddr)
	if err := server.CreateDatabaseStructure(); err != nil {
		logger.Fatalf("Failed to created required database structure: %v\n", err)
	}

	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	db.Close()
	logger.Println("Server stopped")
}
