package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct untuk menyimpan semua konfigurasi aplikasi
type Config struct {
	AppEnv   string
	Port     string
	Database DatabaseConfig
	Email    EmailConfig
}

type DatabaseConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	DBName      string `json:"dbname"`
	SSLMode     string `json:"sslmode"`
	TimeZone    string `json:"timezone"`
	DatabaseURL string
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	UploadDir     string `json:"upload_dir"`
	ProcessedDir  string `json:"processed_dir"`
	ExportDir     string `json:"export_dir"`
	MaxUploadSize int64  `json:"max_upload_size"`
}

type JWTConfig struct {
	Secret     string
	Expiration int // in hours
}

type EmailConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	SenderName   string `mapstructure:"SMTP_SENDER_NAME"`
	AuthEmail    string `mapstructure:"SMTP_AUTH_EMAIL"`
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
}

// LoadConfig memuat konfigurasi dari file .env atau environment variables.
func LoadConfig() (*Config, error) {
	// Memuat file .env. Jika tidak ada atau gagal, lanjut tanpa error (mungkin dari env vars)
	err := godotenv.Load()
	if err != nil {
		log.Println("Tidak ada file .env yang ditemukan, memuat dari environment variables.")
	}

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("APP_PORT"),
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Email: EmailConfig{
			Host:       os.Getenv("SMTP_HOST"),
			Port:       587,
			SenderName: os.Getenv("SMTP_SENDER_NAME"),
		},
	}

	// Bangun URL database dari komponen-komponen atau gunakan DATABASE_URL jika sudah ada
	if os.Getenv("DATABASE_URL") != "" {
		cfg.Database.DatabaseURL = os.Getenv("DATABASE_URL")
	} else {
		// Pastikan semua variabel DB_ diperlukan jika DATABASE_URL tidak diset
		cfg.Database.DatabaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	// Atur port default jika tidak diset di environment
	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	// Log konfigurasi (opsional, untuk debugging)
	log.Printf("Aplikasi berjalan di lingkungan: %s, Port: %s", cfg.AppEnv, cfg.Port)
	return cfg, nil
}
