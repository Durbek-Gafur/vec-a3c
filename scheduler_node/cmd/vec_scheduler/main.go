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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	mysqlStore, err := initMySQLStore(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL store: %v", err)
	}

	if err := populateDatabase(ctx, mysqlStore); err != nil {
		log.Fatalf("Failed to populate database: %v", err)
	}

	server := startServer(mysqlStore)
	waitForShutdown(server)
}

func initMySQLStore(cfg *Config) (*store.MySQLStore, error) {
	dsn := cfg.MySQLUser + ":" + cfg.MySQLPassword + "@tcp(" + cfg.MySQLHost + ":" + cfg.MySQLPort + ")/" + cfg.MySQLDatabase + "?parseTime=true"
	return store.NewMySQLStore(dsn)
}

func populateDatabase(ctx context.Context, store *store.MySQLStore) error {
	seeder := seeder.NewDBSeeder(store, store, &seeder.ActualURLProvider{})

	if err := seeder.PopulateVENInfo(); err != nil {
		return err
	}

	if err := seeder.PopulateWorkflows(); err != nil {
		return err
	}

	// Start the UpdateQueueSizePeriodically function in a separate goroutine.
	go seeder.UpdateQueueSizePeriodically(ctx)

	return nil
}

func startServer(store *store.MySQLStore) *http.Server {
	handler := api.NewHandler(store, store)

	router := mux.NewRouter()
	router.HandleFunc("/", handler.ShowTables).Methods("GET")
	router.HandleFunc("/index", handler.ShowTables).Methods("GET")
	router.HandleFunc("/workflow", handler.UpdateWorkflowByName).Methods("POST")

	server := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	// Start the server in a separate goroutine.
	go func() {
		log.Println("Started the app on port 8090")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	return server
}

func waitForShutdown(server *http.Server) {
	// Set up a channel to listen for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	// Initiate a graceful shutdown.
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down the server gracefully: %v", err)
	}
	log.Println("Server stopped")
}
