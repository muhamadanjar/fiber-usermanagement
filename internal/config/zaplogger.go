package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig represents logging configuration
type LogConfig struct {
	Level      *string `json:"level" mapstructure:"level"`             // debug, info, warn, error
	Format     *string `json:"format" mapstructure:"format"`           // json, console
	Output     *string `json:"output" mapstructure:"output"`           // stdout, stderr, file path
	MaxSize    *int    `json:"max_size" mapstructure:"max_size"`       // Max size in MB
	MaxBackups *int    `json:"max_backups" mapstructure:"max_backups"` // Max backup files
	MaxAge     *int    `json:"max_age" mapstructure:"max_age"`         // Max age in days
	Compress   *bool   `json:"compress" mapstructure:"compress"`       // Compress old files
}

// Add LogConfig to your main Config struct
// Update your Config struct to include:
// Log LogConfig `mapstructure:"log"`

// ZapLoggerManager manages zap logger configuration
type ZapLoggerManager struct {
	config *Config
	logger *zap.Logger
}

// NewZapLogger creates a new zap logger with configuration support
func NewZapLogger(config *Config) (*zap.Logger, error) {
	manager := &ZapLoggerManager{
		config: config,
	}

	logger, err := manager.createLogger()
	if err != nil {
		return nil, err
	}

	manager.logger = logger
	return logger, nil
}

// NewZapLoggerWithDefaults creates a logger with default configuration if config is nil
func NewZapLoggerWithDefaults() (*zap.Logger, error) {
	// Create default config if none provided
	config := &Config{}
	setDefaultLogConfig(config)

	return NewZapLogger(config)
}

// createLogger creates the actual zap logger based on configuration
func (lm *ZapLoggerManager) createLogger() (*zap.Logger, error) {
	// Set defaults if not configured
	if lm.config.Log.Level == nil || lm.config.Log.Format == nil {
		setDefaultLogConfig(lm.config)
	}

	// Determine if it's development mode
	isDev := lm.config.IsDevelopment()

	// Create encoder config
	encoderConfig := lm.getEncoderConfig(isDev)

	// Create core
	core, err := lm.createCore(encoderConfig)
	if err != nil {
		return nil, err
	}

	// Create logger with options
	options := lm.getLoggerOptions(isDev)
	logger := zap.New(core, options...)

	return logger, nil
}

// getEncoderConfig returns encoder configuration based on environment
func (lm *ZapLoggerManager) getEncoderConfig(isDev bool) zapcore.EncoderConfig {
	if isDev {
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncodeCaller = zapcore.ShortCallerEncoder
		return config
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	return config
}

// createCore creates zapcore.Core based on configuration
func (lm *ZapLoggerManager) createCore(encoderConfig zapcore.EncoderConfig) (zapcore.Core, error) {
	// Get log level
	level := lm.getLogLevel()

	// Get encoder
	encoder := lm.getEncoder(encoderConfig)

	// Get writer syncer
	writeSyncer, err := lm.getWriteSyncer()
	if err != nil {
		return nil, err
	}

	return zapcore.NewCore(encoder, writeSyncer, level), nil
}

// getLogLevel returns the configured log level
func (lm *ZapLoggerManager) getLogLevel() zapcore.Level {
	if lm.config.Log.Level == nil {
		if lm.config.IsDevelopment() {
			return zapcore.DebugLevel
		}
		return zapcore.InfoLevel
	}

	switch *lm.config.Log.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoder returns the configured encoder
func (lm *ZapLoggerManager) getEncoder(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	format := "json"
	if lm.config.Log.Format != nil {
		format = *lm.config.Log.Format
	}

	// Force console format for development
	if lm.config.IsDevelopment() {
		format = "console"
	}

	switch format {
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// getWriteSyncer returns the configured write syncer
func (lm *ZapLoggerManager) getWriteSyncer() (zapcore.WriteSyncer, error) {
	output := "stdout"
	if lm.config.Log.Output != nil {
		output = *lm.config.Log.Output
	}

	switch output {
	case "stdout":
		return zapcore.AddSync(os.Stdout), nil
	case "stderr":
		return zapcore.AddSync(os.Stderr), nil
	default:
		// Assume it's a file path
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", output, err)
		}
		return zapcore.AddSync(file), nil
	}
}

// getLoggerOptions returns logger options based on environment
func (lm *ZapLoggerManager) getLoggerOptions(isDev bool) []zap.Option {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if isDev {
		options = append(options, zap.Development())
	}

	return options
}

// setDefaultLogConfig sets default logging configuration
func setDefaultLogConfig(config *Config) {
	if config.Log.Level == nil {
		level := "info"
		if config.IsDevelopment() {
			level = "debug"
		}
		config.Log.Level = &level
	}

	if config.Log.Format == nil {
		format := "json"
		if config.IsDevelopment() {
			format = "console"
		}
		config.Log.Format = &format
	}

	if config.Log.Output == nil {
		output := "stdout"
		config.Log.Output = &output
	}

	if config.Log.MaxSize == nil {
		maxSize := 100 // 100MB
		config.Log.MaxSize = &maxSize
	}

	if config.Log.MaxBackups == nil {
		maxBackups := 3
		config.Log.MaxBackups = &maxBackups
	}

	if config.Log.MaxAge == nil {
		maxAge := 28 // 28 days
		config.Log.MaxAge = &maxAge
	}

	if config.Log.Compress == nil {
		compress := true
		config.Log.Compress = &compress
	}
}

// LoggerWithContext creates a logger with additional context fields
func LoggerWithContext(logger *zap.Logger, fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// GetGlobalLogger returns a global logger instance (singleton pattern)
var globalLogger *zap.Logger

func InitGlobalLogger(config *Config) error {
	logger, err := NewZapLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func GetGlobalLogger() *zap.Logger {
	if globalLogger == nil {
		// Fallback to default logger if not initialized
		logger, err := NewZapLoggerWithDefaults()
		if err != nil {
			// Last resort: use nop logger
			return zap.NewNop()
		}
		globalLogger = logger
	}
	return globalLogger
}

// SyncGlobalLogger syncs the global logger
func SyncGlobalLogger() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Helper functions for common logging patterns

// LogError logs an error with context
func LogError(logger *zap.Logger, msg string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	logger.Error(msg, allFields...)
}

// LogErrorf logs an error with formatted message
func LogErrorf(logger *zap.Logger, format string, err error, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...), zap.Error(err))
}

// LogWithRequestID logs with request ID context
func LogWithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	return logger.With(zap.String("request_id", requestID))
}

// LogWithUserID logs with user ID context
func LogWithUserID(logger *zap.Logger, userID string) *zap.Logger {
	return logger.With(zap.String("user_id", userID))
}
