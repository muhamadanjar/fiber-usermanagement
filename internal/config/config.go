package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents the main configuration structure
type Config struct {
	AppEnv   *string        `mapstructure:"app_env"`
	Port     *string        `mapstructure:"port"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Email    EmailConfig    `mapstructure:"email"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Log      LogConfig      `mapstructure:"log"` // Add LogConfig here
	Redis    RedisConfig    `mapstructure:"redis"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host        *string `json:"host" mapstructure:"host"`
	Port        *string `json:"port" mapstructure:"port"`
	User        *string `json:"user" mapstructure:"user"`
	Password    *string `json:"password" mapstructure:"password"`
	DBName      *string `json:"dbname" mapstructure:"dbname"`
	SSLMode     *string `json:"sslmode" mapstructure:"sslmode"`
	TimeZone    *string `json:"timezone" mapstructure:"timezone"`
	DatabaseURL *string `mapstructure:"database_url"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	UploadDir     *string `json:"upload_dir" mapstructure:"upload_dir"`
	ProcessedDir  *string `json:"processed_dir" mapstructure:"processed_dir"`
	ExportDir     *string `json:"export_dir" mapstructure:"export_dir"`
	MaxUploadSize *int64  `json:"max_upload_size" mapstructure:"max_upload_size"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret     *string `mapstructure:"secret"`
	Expiration *int    `mapstructure:"expiration"` // in hours
}

// EmailConfig represents email configuration
type EmailConfig struct {
	Host         *string `mapstructure:"host"`
	Port         *int    `mapstructure:"port"`
	SenderName   *string `mapstructure:"sender_name"`
	AuthEmail    *string `mapstructure:"auth_email"`
	AuthPassword *string `mapstructure:"auth_password"`
}

type RabbitMQConfig struct {
	URL *string `mapstructure:"url"`
}

type RedisConfig struct {
	Host     *string `json:"host" mapstructure:"host"`
	Port     *int    `json:"port" mapstructure:"port"`
	Password *string `json:"password" mapstructure:"password"`
	DB       *int    `json:"db" mapstructure:"db"`
}

// ConfigManager handles configuration loading and management
type ConfigManager struct {
	viper  *viper.Viper
	config *Config
}

// NewConfig creates a new config instance with support for JSON, ENV files, and environment variables
func NewConfig() (*Config, error) {
	manager := &ConfigManager{
		viper: viper.New(),
	}

	// Load .env file if exists (optional, won't panic if not found)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or unable to load: %v", err)
	}

	// Setup and load configuration
	manager.setupViper()
	if err := manager.loadConfig(); err != nil {
		return nil, err
	}

	return manager.config, nil
}

// MustLoadConfig loads config and panics on error (for backward compatibility)
func MustLoadConfig() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return config
}

// setupViper configures viper settings
func (cm *ConfigManager) setupViper() {
	// Set config file settings
	cm.viper.SetConfigName("config")
	cm.viper.SetConfigType("json")
	cm.viper.AddConfigPath("./../")
	cm.viper.AddConfigPath("./")
	cm.viper.AddConfigPath("./config/")

	// Enable environment variable support
	cm.viper.AutomaticEnv()

	// Set environment variable prefix
	cm.viper.SetEnvPrefix("APP")

	// Replace characters in environment variable names
	cm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Set default values
	cm.setDefaults()
}

// setDefaults sets default configuration values
func (cm *ConfigManager) setDefaults() {
	// App defaults
	cm.viper.SetDefault("app_env", "development")
	cm.viper.SetDefault("port", "8080")

	// Database defaults
	cm.viper.SetDefault("database.host", "localhost")
	cm.viper.SetDefault("database.port", "5432")
	cm.viper.SetDefault("database.user", "postgres")
	cm.viper.SetDefault("database.password", "")
	cm.viper.SetDefault("database.dbname", "myapp")
	cm.viper.SetDefault("database.sslmode", "disable")
	cm.viper.SetDefault("database.timezone", "UTC")

	// Storage defaults
	cm.viper.SetDefault("storage.upload_dir", "./uploads")
	cm.viper.SetDefault("storage.processed_dir", "./processed")
	cm.viper.SetDefault("storage.export_dir", "./exports")
	cm.viper.SetDefault("storage.max_upload_size", int64(10*1024*1024)) // 10MB

	// JWT defaults
	cm.viper.SetDefault("jwt.secret", "your-secret-key")
	cm.viper.SetDefault("jwt.expiration", 24) // 24 hours

	// Email defaults
	cm.viper.SetDefault("email.host", "localhost")
	cm.viper.SetDefault("email.port", 587)
	cm.viper.SetDefault("email.sender_name", "MyApp")
	cm.viper.SetDefault("email.auth_email", "")
	cm.viper.SetDefault("email.auth_password", "")

	// Redis defaults
	cm.viper.SetDefault("redis.host", "localhost")
	cm.viper.SetDefault("redis.port", 6379)
	cm.viper.SetDefault("redis.password", "")
	cm.viper.SetDefault("redis.db", 0)
}

// loadConfig loads configuration from various sources and unmarshals to struct
func (cm *ConfigManager) loadConfig() error {
	// Try to read config file (optional)
	if err := cm.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Warning: Config file not found, using defaults and environment variables")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", cm.viper.ConfigFileUsed())
	}

	// Unmarshal config to struct
	config := &Config{}
	if err := cm.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Post-process configuration
	cm.postProcessConfig(config)

	cm.config = config
	return nil
}

// postProcessConfig performs post-processing on loaded configuration
func (cm *ConfigManager) postProcessConfig(config *Config) {
	// Generate database URL if not provided
	if config.Database.DatabaseURL == nil || *config.Database.DatabaseURL == "" {
		if config.Database.Host != nil && config.Database.Port != nil &&
			config.Database.User != nil && config.Database.DBName != nil {
			dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s",
				getStringValue(config.Database.Host),
				getStringValue(config.Database.Port),
				getStringValue(config.Database.User),
				getStringValue(config.Database.Password),
				getStringValue(config.Database.DBName),
				getStringValue(config.Database.SSLMode),
				getStringValue(config.Database.TimeZone),
			)
			config.Database.DatabaseURL = &dbURL
		}
	}
}

// Helper methods for the Config struct

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	if c.Database.DatabaseURL != nil {
		return *c.Database.DatabaseURL
	}
	return ""
}

// GetServerAddress returns formatted server address
func (c *Config) GetServerAddress() string {
	port := "8080"
	if c.Port != nil {
		port = *c.Port
	}
	return ":" + port
}

// IsProduction returns whether the app is running in production
func (c *Config) IsProduction() bool {
	return c.AppEnv != nil && *c.AppEnv == "production"
}

// IsDevelopment returns whether the app is running in development
func (c *Config) IsDevelopment() bool {
	return c.AppEnv != nil && *c.AppEnv == "development"
}

// ValidateConfig validates required configuration values
func (c *Config) ValidateConfig() error {
	// Validate required fields
	if c.Database.Host == nil || *c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.DBName == nil || *c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if c.JWT.Secret == nil || *c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// PrintConfig prints current configuration (for debugging)
func (c *Config) PrintConfig() {
	fmt.Println("Current Configuration:")
	fmt.Printf("  App Environment: %s\n", getStringValue(c.AppEnv))
	fmt.Printf("  Port: %s\n", getStringValue(c.Port))

	fmt.Println("  Database:")
	fmt.Printf("    Host: %s\n", getStringValue(c.Database.Host))
	fmt.Printf("    Port: %s\n", getStringValue(c.Database.Port))
	fmt.Printf("    User: %s\n", getStringValue(c.Database.User))
	fmt.Printf("    Password: ****\n")
	fmt.Printf("    DBName: %s\n", getStringValue(c.Database.DBName))
	fmt.Printf("    SSLMode: %s\n", getStringValue(c.Database.SSLMode))
	fmt.Printf("    TimeZone: %s\n", getStringValue(c.Database.TimeZone))

	fmt.Println("  Storage:")
	fmt.Printf("    Upload Dir: %s\n", getStringValue(c.Storage.UploadDir))
	fmt.Printf("    Processed Dir: %s\n", getStringValue(c.Storage.ProcessedDir))
	fmt.Printf("    Export Dir: %s\n", getStringValue(c.Storage.ExportDir))
	fmt.Printf("    Max Upload Size: %d\n", getInt64Value(c.Storage.MaxUploadSize))

	fmt.Println("  JWT:")
	fmt.Printf("    Secret: ****\n")
	fmt.Printf("    Expiration: %d hours\n", getIntValue(c.JWT.Expiration))

	fmt.Println("  Email:")
	fmt.Printf("    Host: %s\n", getStringValue(c.Email.Host))
	fmt.Printf("    Port: %d\n", getIntValue(c.Email.Port))
	fmt.Printf("    Sender Name: %s\n", getStringValue(c.Email.SenderName))
	fmt.Printf("    Auth Email: %s\n", getStringValue(c.Email.AuthEmail))
	fmt.Printf("    Auth Password: ****\n")
}

// Helper functions to safely get values from pointers
func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func getIntValue(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func getInt64Value(ptr *int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
