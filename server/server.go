package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ApplicationServer struct {
	Logger    *log.Logger
	Db        *pgxpool.Pool
	webserver *http.Server
	router    *gin.Engine
}

// NewApplicationServer creates a new server with the given configuration.
// listenAddr example: ":5000"
func NewApplicationServer(logger *log.Logger, db *pgxpool.Pool, listenAddr string) ApplicationServer {
	// create router
	router := gin.Default()
	router.Use(gin.Logger())

	// create application server
	server := ApplicationServer{
		Logger: logger,
		webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		Db:     db,
		router: router,
	}

	// configure routes
	router.GET("/", welcome)

	router.GET("/user", server.getUsers)
	router.GET("/user/:id", server.getUser)
	router.DELETE("/user/:id", server.deleteUser)
	router.POST("/user", server.addUser)

	//router.HandleFunc("/", logCall(logger, welcome))
	//router.HandleFunc("/user", logCall(logger, server.user))

	return server
}

func (srv ApplicationServer) GracefullShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	server := srv.webserver
	logger := srv.Logger

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
