package main

import (
	"auth-service/config"
	"auth-service/src/infrastructure/database"
	"auth-service/src/registry"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// $Env:GOOS = "linux"; $Env:GOARCH = "386"; go build
func main() {

	err := config.LoadConfig()

	if err != nil {
		log.Fatal("error loading environment variables:", err)
		os.Exit(1)
	}

	port := config.C.Server.Address

	db, err := database.NewDB()

	if err != nil {
		log.Fatal("Error when starting database:", err)
	}

	reg := registry.NewRegistry(db)
	router, err := reg.StartHandlerInterface()

	if err != nil {
		log.Fatal("Error when starting app:", err)
	}

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	log.Printf("listen: %s\n", port)

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
