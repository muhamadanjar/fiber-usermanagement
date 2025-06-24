package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct untuk menyimpan semua konfigurasi aplikasi
type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
}

// LoadConfig memuat konfigurasi dari file .env atau environment variables.
func LoadConfig() *Config {
	// Memuat file .env. Jika tidak ada atau gagal, lanjut tanpa error (mungkin dari env vars)
	err := godotenv.Load()
	if err != nil {
		log.Println("Tidak ada file .env yang ditemukan, memuat dari environment variables.")
	}

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		Port:   os.Getenv("APP_PORT"),
	}

	// Bangun URL database dari komponen-komponen atau gunakan DATABASE_URL jika sudah ada
	if os.Getenv("DATABASE_URL") != "" {
		cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	} else {
		// Pastikan semua variabel DB_ diperlukan jika DATABASE_URL tidak diset
		cfg.DatabaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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
	return cfg
}
