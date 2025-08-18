package main

import (
	"context"
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/container"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	appContainer, err := config.NewAppContainer()
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}

	// Ensure cleanup on exit
	defer func() {
		if err := appContainer.Close(); err != nil {
			appContainer.Logger.Error("Error during application shutdown", zap.Error(err))
		}
	}()

	businessContainer, err := container.NewContainer(appContainer)
	if err != nil {
		appContainer.Logger.Fatal("Failed to initialize business container", zap.Error(err))
	}

	// Setup routes
	businessContainer.SetupRoutes()

	// Start server in a goroutine
	serverAddr := appContainer.GetServerAddress()
	go func() {
		appContainer.Logger.Info("Starting server", zap.String("address", serverAddr))
		if err := appContainer.App.Listen(serverAddr); err != nil {
			appContainer.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal
	<-c
	appContainer.Logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := appContainer.App.ShutdownWithContext(ctx); err != nil {
		appContainer.Logger.Error("Server forced to shutdown", zap.Error(err))
	}

	appContainer.Logger.Info("Server exited")
}
