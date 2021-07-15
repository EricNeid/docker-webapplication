package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ApplicationServer struct {
	logger    *log.Logger
	db        *pgxpool.Pool
	webserver *http.Server
	router    *gin.Engine
}

// NewApplicationServer creates a new server with the given configuration.
// listenAddr example: ":5000"
func NewApplicationServer(db *pgxpool.Pool, listenAddr string) ApplicationServer {
	// create logger
	logger := log.New(os.Stdout, "server", log.LstdFlags)

	// create router
	router := gin.Default()
	router.Use(gin.Logger())

	// create application server
	server := ApplicationServer{
		logger: logger,
		webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		db:     db,
		router: router,
	}

	// configure routes
	router.GET("/", welcome)

	// user crud
	router.GET("/users", server.getUsers)
	router.GET("/users/:id", server.getUser)
	router.DELETE("/users/:id", server.deleteUser)
	router.POST("/users", server.addUser)

	// vehicle state crud
	router.GET("/vehicleStates", server.getVehicleStates)
	router.GET("/vehicleStates/:id", server.getVehicleState)
	router.DELETE("/vehicleStates/:id", server.deleteVehicleState)
	router.POST("/vehicleStates", server.addVehicleState)

	return server
}

func (srv ApplicationServer) SetLogWriter(out io.Writer) {
	srv.logger.SetOutput(out)
	gin.DefaultWriter = out
}

func (srv ApplicationServer) CreateDatabaseStructure() error {
	logger := srv.logger
	db := srv.db
	err := createTableVehicleState(logger, db)
	if err != nil {
		return err
	}
	err = createTableUsers(logger, db)
	return err
}

func (srv ApplicationServer) GracefullShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	server := srv.webserver
	logger := srv.logger

	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	close(done)
}

func (srv ApplicationServer) ListenAndServe() error {
	return srv.webserver.ListenAndServe()
}
