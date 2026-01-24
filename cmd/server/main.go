package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/simple-wallet/internal/config"
	"github.com/simple-wallet/internal/handler"
	"github.com/simple-wallet/internal/repository"
	"github.com/simple-wallet/internal/service"
)

func main() {
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Printf("Warning: failed to load config.env: %v, using environment variables", err)
		cfg, err = config.Load("")
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var pool *repository.DBPool
	for i := 0; i < 30; i++ {
		pool, err = repository.NewPostgresPool(ctx, &cfg.Database)
		if err == nil {
			break
		}
		log.Printf("Waiting for database... attempt %d/30: %v", i+1, err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after 30 attempts: %v", err)
	}
	defer pool.Close()

	if err := repository.RunMigrations(&cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	walletRepo := repository.NewWalletRepository(pool)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	walletHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
