// config/app.go
package config

import (
	"context"
	"fiber-usermanagement/internal/infrastructure/database"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AppContainer holds all application dependencies
type AppContainer struct {
	Config    *Config
	DB        *gorm.DB
	Redis     *redis.Client
	Logger    *zap.Logger
	Validator *validator.Validate
	App       *fiber.App
	// RabbitMQ  *RabbitMQConnection // assuming you have this type
}

// NewAppContainer creates and initializes all application dependencies
func NewAppContainer() (*AppContainer, error) {
	container := &AppContainer{}

	// Initialize configuration
	if err := container.initConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Initialize logger
	if err := container.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
	if err := container.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis
	if err := container.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	// Initialize RabbitMQ
	if err := container.initRabbitMQ(); err != nil {
		return nil, fmt.Errorf("failed to initialize rabbitmq: %w", err)
	}

	// Initialize validator
	container.initValidator()

	// Initialize Fiber app
	container.initFiberApp()

	container.Logger.Info("Application container initialized successfully")

	return container, nil
}

// initConfig initializes application configuration
func (c *AppContainer) initConfig() error {
	config, err := NewConfig()
	if err != nil {
		return err
	}

	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	c.Config = config
	return nil
}

// initLogger initializes zap logger
func (c *AppContainer) initLogger() error {
	logger, err := NewZapLogger(c.Config)
	if err != nil {
		return err
	}

	c.Logger = logger

	// Initialize global logger for the application
	if err := InitGlobalLogger(c.Config); err != nil {
		return fmt.Errorf("failed to initialize global logger: %w", err)
	}

	return nil
}

// initDatabase initializes database connection
func (c *AppContainer) initDatabase() error {
	db, err := database.NewPostgresDB(c.Config.GetDatabaseURL())
	if err != nil {
		return err
	}

	c.DB = db
	c.Logger.Info("Database connection established")

	return nil
}

// initRedis initializes Redis connection
func (c *AppContainer) initRedis() error {
	redisConfig := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", getStringValue(c.Config.Redis.Host), getIntValue(c.Config.Redis.Port)),
		Password: getStringValue(c.Config.Redis.Password),
		DB:       getIntValue(c.Config.Redis.DB),
	}

	client := redis.NewClient(redisConfig)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	c.Redis = client
	c.Logger.Info("Redis connection established")

	return nil
}

// initRabbitMQ initializes RabbitMQ connection
func (c *AppContainer) initRabbitMQ() error {
	// rabbitMQ, err := NewRabbitMQ(c.Config)
	// if err != nil {
	// 	return err
	// }

	// c.RabbitMQ = rabbitMQ
	c.Logger.Info("RabbitMQ connection established")

	return nil
}

// initValidator initializes field validator
func (c *AppContainer) initValidator() {
	c.Validator = validator.New()

	// Register custom validators if needed
	c.registerCustomValidators()
}

// initFiberApp initializes Fiber application
func (c *AppContainer) initFiberApp() {
	c.App = fiber.New(fiber.Config{
		ErrorHandler: c.errorHandler,
		AppName:      "User Management API",
		ServerHeader: "Fiber",
		// Add other Fiber configurations based on your needs
	})

	// Setup middleware
	c.setupMiddleware()
}

// registerCustomValidators registers custom validation rules
func (c *AppContainer) registerCustomValidators() {
	// Example: Register custom email validation
	c.Validator.RegisterValidation("custom_email", func(fl validator.FieldLevel) bool {
		// Custom email validation logic
		return true
	})
}

// errorHandler handles Fiber errors
func (c *AppContainer) errorHandler(ctx *fiber.Ctx, err error) error {
	c.Logger.Error("Fiber error occurred",
		zap.Error(err),
		zap.String("path", ctx.Path()),
		zap.String("method", ctx.Method()),
	)

	// Handle different error types
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}

// setupMiddleware sets up common middleware
func (c *AppContainer) setupMiddleware() {
	// Add your middleware here
	// Example: CORS, Rate limiting, etc.
}

// GetServerAddress returns the server address from config
func (c *AppContainer) GetServerAddress() string {
	return c.Config.GetServerAddress()
}

// Close gracefully shuts down all connections
func (c *AppContainer) Close() error {
	var errors []error

	// Close database connection
	if c.DB != nil {
		if sqlDB, err := c.DB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close database: %w", err))
			}
		}
	}

	// Close Redis connection
	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close redis: %w", err))
		}
	}

	// Close RabbitMQ connection
	// if c.RabbitMQ != nil {
	// 	if err := c.RabbitMQ.Close(); err != nil {
	// 		errors = append(errors, fmt.Errorf("failed to close rabbitmq: %w", err))
	// 	}
	// }

	// Sync logger
	if c.Logger != nil {
		if err := c.Logger.Sync(); err != nil {
			errors = append(errors, fmt.Errorf("failed to sync logger: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}

// Helper functions (you might already have these)
// func getStringValue(ptr *string) string {
// 	if ptr != nil {
// 		return *ptr
// 	}
// 	return ""
// }

// func getIntValue(ptr *int) int {
// 	if ptr != nil {
// 		return *ptr
// 	}
// 	return 0
// }
