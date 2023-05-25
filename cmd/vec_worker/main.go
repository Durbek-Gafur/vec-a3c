package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"vec-node/internal/api"
	"vec-node/internal/store"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	dsn := cfg.MySQLUser + ":" + cfg.MySQLPassword + "@tcp(" + cfg.MySQLHost + ":" + cfg.MySQLPort + ")/" + cfg.MySQLDatabase + "?parseTime=true"
	mysqlStore, err := store.NewMySQLStore(dsn)
	if err != nil {
		log.Fatal("Failed to initialize MySQL store:", err)
	}

	handler := api.NewHandler(mysqlStore, mysqlStore)

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
		log.Printf("Started the app")
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
