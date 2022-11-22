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

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	listenAddr string = ":5000"
	dbHost     string = "localhost"
	dbPort     int    = 5432
	dbUser     string = "postgres"
	dbPass     string = "postgres"
	dbName     string = "localdb"

	logFile string = ""
)

func init() {
	readConfigFromEnvironment()
	readConfigFromCli()
}

func main() {
	// prepare logging and gin
	var logOut io.Writer
	if logFile != "" {
		logOut = io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    500, // megabytes
				MaxBackups: 3,
				MaxAge:     28, //days
			},
		)
	} else {
		logOut = os.Stdout
	}
	gin.DefaultWriter = logOut
	log.SetOutput(logOut)
	log.SetPrefix("[APP] ")

	// prepare shutdown channel
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// create db pool
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	log.Printf("Connecting to db using: %s", dbURL)
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Could not create database pool: %v\n", err)
	}

	// create server
	gin.SetMode(gin.ReleaseMode)
	server := server.NewApplicationServer(db, listenAddr)
	go server.GracefullShutdown(quit, done)

	log.Println("Creating database structure", listenAddr)
	if err := server.CreateDatabaseStructure(); err != nil {
		log.Fatalf("Failed to created required database structure: %v\n", err)
	}

	log.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	db.Close()
	log.Println("Server stopped")
}

func readConfigFromEnvironment() {
	if value, isSet := os.LookupEnv("LISTEN_ADDR"); isSet {
		listenAddr = value
	}

	if value, isSet := os.LookupEnv("DB_HOST"); isSet {
		dbHost = value
	}

	if value, isSet := os.LookupEnv("DB_PORT"); isSet {
		dbPort, _ = strconv.Atoi(value)
	}

	if value, isSet := os.LookupEnv("DB_USER"); isSet {
		dbUser = value
	}

	if value, isSet := os.LookupEnv("DB_PASS"); isSet {
		dbPass = value
	}

	if value, isSet := os.LookupEnv("DB_NAME"); isSet {
		dbName = value
	}

	if value, isSet := os.LookupEnv("LOG_FILE"); isSet {
		logFile = value
	}
}

func readConfigFromCli() {
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.StringVar(&dbHost, "db-host", dbHost, "database host address")
	flag.IntVar(&dbPort, "db-port", dbPort, "database port")
	flag.StringVar(&dbUser, "db-user", dbUser, "database user credential")
	flag.StringVar(&dbPass, "db-pass", dbPass, "database user password")
	flag.StringVar(&dbName, "db-name", dbName, "database name")
	flag.StringVar(&logFile, "log-file", logFile, "Optional: write log to this file")

	flag.Parse()
}
