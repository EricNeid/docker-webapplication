package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const timeFormat string = "Mon Jan 2 15:04:05 2006"

type ApplicationServer struct {
	Webserver *http.Server
	Logger    *log.Logger
	db        *pgxpool.Pool
}

// NewApplicationServer creates a new server with the given configuration.
// listenAddr example: ":5000"
func NewApplicationServer(logger *log.Logger, db *pgxpool.Pool, listenAddr string) ApplicationServer {
	// create router
	router := http.NewServeMux()

	// create application server
	server := ApplicationServer{
		Logger: logger,
		Webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ErrorLog:     logger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		db: db,
	}

	// configure routes
	router.HandleFunc("/", logCall(logger, welcome))
	router.HandleFunc("/user", logCall(logger, server.user))

	return server
}

func (srv ApplicationServer) GracefullShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	server := srv.Webserver
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
	return srv.Webserver.ListenAndServe()
}

func logCall(logger *log.Logger, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		logger.Printf("%s - %s\n", timestamp.Format(timeFormat), r.URL.Path)
		handler(w, r)
	}
}
