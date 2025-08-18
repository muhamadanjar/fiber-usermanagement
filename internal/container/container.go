// container/container.go
package container

import (
	"fiber-usermanagement/internal/api/handlers"
	"fiber-usermanagement/internal/api/routes"
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/domain/repositories"
	"fiber-usermanagement/internal/infrastructure/persistence"
	"fiber-usermanagement/internal/usecase/interactors"
	"fmt"
)

// BusinessContainer holds all business logic dependencies
type BusinessContainer struct {
	appContainer *config.AppContainer

	// Repositories
	userRepo repositories.UserRepository

	// Interactors/Use Cases
	userInteractor *interactors.UserInteractor

	// Handlers
	userHandler *handlers.UserHandler
}

// NewContainer creates a new business container with all dependencies
func NewContainer(appContainer *config.AppContainer) (*BusinessContainer, error) {
	container := &BusinessContainer{
		appContainer: appContainer,
	}

	// Initialize repositories
	if err := container.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	// Initialize interactors
	if err := container.initInteractors(); err != nil {
		return nil, fmt.Errorf("failed to initialize interactors: %w", err)
	}

	// Initialize handlers
	if err := container.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	// Run database migrations
	if err := container.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	appContainer.Logger.Info("Business container initialized successfully")

	return container, nil
}

// initRepositories initializes all repository implementations
func (c *BusinessContainer) initRepositories() error {
	c.userRepo = persistence.NewUserRepository(c.appContainer.DB)

	c.appContainer.Logger.Info("Repositories initialized")
	return nil
}

// initInteractors initializes all use case interactors
func (c *BusinessContainer) initInteractors() error {
	c.userInteractor = interactors.NewUserInteractor(c.userRepo)

	c.appContainer.Logger.Info("Interactors initialized")
	return nil
}

// initHandlers initializes all HTTP handlers
func (c *BusinessContainer) initHandlers() error {
	c.userHandler = handlers.NewUserHandler(c.userInteractor)

	c.appContainer.Logger.Info("Handlers initialized")
	return nil
}

// runMigrations runs database migrations
func (c *BusinessContainer) runMigrations() error {
	entities := []interface{}{
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
	}

	for _, entity := range entities {
		if err := c.appContainer.DB.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", entity, err)
		}
	}

	c.appContainer.Logger.Info("Database migrations completed successfully")
	return nil
}

// SetupRoutes configures all application routes
func (c *BusinessContainer) SetupRoutes() {
	routeConfig := &routes.RouteConfig{
		App: c.appContainer.App,
		// Logger:      c.appContainer.Logger,
		UserHandler: c.userHandler,
		// Add other handlers as needed
	}

	routeConfig.Setup()
	c.appContainer.Logger.Info("Routes configured successfully")
}

// GetAppContainer returns the application container
func (c *BusinessContainer) GetAppContainer() *config.AppContainer {
	return c.appContainer
}

// Health check methods
func (c *BusinessContainer) HealthCheck() error {
	// Check database connection
	if sqlDB, err := c.appContainer.DB.DB(); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	} else {
		if err := sqlDB.Ping(); err != nil {
			return fmt.Errorf("database ping error: %w", err)
		}
	}

	// Check Redis connection
	// if err := c.appContainer.Redis.Ping(c.appContainer.Redis.Context()).Err(); err != nil {
	// 	return fmt.Errorf("redis connection error: %w", err)
	// }

	// Add other health checks as needed

	return nil
}
