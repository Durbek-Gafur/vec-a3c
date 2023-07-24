package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"vec-node/internal/api"
	"vec-node/internal/rspec"
	"vec-node/internal/store"
	"vec-node/internal/workflow"
)

type MultiWriter struct {
	Writers []io.Writer
}

func (mw MultiWriter) Write(p []byte) (n int, err error) {
	for _, writer := range mw.Writers {
		writer.Write(p)
	}
	return len(p), nil
}

func main() {
	logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	multiWriter := MultiWriter{
		Writers: []io.Writer{
			os.Stdout,
			logFile,
		},
	}
	log.SetOutput(multiWriter)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	dsn := cfg.MySQLUser + ":" + cfg.MySQLPassword + "@tcp(" + cfg.MySQLHost + ":" + cfg.MySQLPort + ")/" + cfg.MySQLDatabase + "?parseTime=true"
	mysqlStore, err := store.NewMySQLStore(dsn)
	if err != nil {
		log.Fatal("Failed to initialize MySQL store:", err)
	}
	wfs := workflow.NewService(mysqlStore, logFile)
	if err != nil {
		log.Fatal("Failed to initialize Workflow Provider:", err)
	}
	rspec_provider := rspec.NewService()
	handler := api.NewHandler(mysqlStore, mysqlStore, rspec_provider, mysqlStore, wfs)

	router := mux.NewRouter()
	// queue
	router.HandleFunc("/queue-size", handler.GetQueueSize).Methods("GET")
	router.HandleFunc("/queue-size", handler.SetQueueSize).Methods("POST")
	router.HandleFunc("/queue-size", handler.UpdateQueueSize).Methods("PUT")

	// workflow
	// TODO pagination for workflows
	router.HandleFunc("/workflow", handler.SaveWorkflow).Methods("POST")
	router.HandleFunc("/workflow/{id:[0-9]+}", handler.GetWorkflowByID).Methods("GET")
	router.HandleFunc("/workflows", handler.GetWorkflows).Methods("GET")

	// rspec
	router.HandleFunc("/rspec", handler.GetRspec).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a separate goroutine.
	go func() {
		log.Printf("Started the app.")
		log.Printf("Server running at port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start the server:", err)
		}
	}()

	// Set up a channel to listen for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	// Initiate a graceful shutdown.
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shut down the server gracefully:", err)
	}
	log.Println("Server stopped")
}
