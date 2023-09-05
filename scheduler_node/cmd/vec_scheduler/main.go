package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scheduler-node/internal/api"
	"scheduler-node/internal/seeder"
	"scheduler-node/internal/store"

	"github.com/gorilla/mux"
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

	//populate table
	if err := seeder.PopulateVENInfo(mysqlStore.GetDB(), &seeder.ActualURLProvider{}); err != nil {
		log.Fatal("Failed to populate VEN table", err)
	}

	if err := seeder.PopulateWorkflows(mysqlStore.GetDB()); err != nil {
		log.Fatal("Failed to populate Workflow table", err)
	}

	// Start the UpdateQueueSizePeriodically function in a separate goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go seeder.UpdateQueueSizePeriodically(ctx, mysqlStore.GetDB(), &seeder.ActualURLProvider{})

	router := mux.NewRouter()
	// queue
	router.HandleFunc("/", handler.ShowTables).Methods("GET")
	router.HandleFunc("/index", handler.ShowTables).Methods("GET")
	router.HandleFunc("/workflow", handler.UpdateWorkflowByName).Methods("POST")

	server := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	// Start the server in a separate goroutine.
	go func() {
		log.Printf("Started the app")
		log.Printf("Server running at port 8090")
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shut down the server gracefully:", err)
	}
	log.Println("Server stopped")
}
