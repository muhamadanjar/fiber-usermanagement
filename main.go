package main

import (
	"fiber-usermanagement/api/routes"
	"fiber-usermanagement/config"
	"fiber-usermanagement/domain/entities"
	"fiber-usermanagement/infrastructure/database"
	"fiber-usermanagement/infrastructure/persistence"
	"fiber-usermanagement/usecase/interactors"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	// Inisialisasi koneksi database
	db, err := database.NewDBConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Migrasi database
	// Tambahkan semua model entitas Anda di sini agar GORM dapat membuatnya di database.
	database.Migrate(db, &entities.User{})

	// Inisialisasi Fiber app
	app := fiber.New()

	// ==== Setup Clean Architecture Layers ====

	// Infrastructure Layer (Persistence)
	// Inisialisasi implementasi repository
	userRepo := persistence.NewUserRepository(db)
	// productRepo := persistence.NewProductRepository(db) // Contoh repository lain

	// Usecase Layer (Interactors)
	// Inisialisasi use case dengan repository yang telah diinisialisasi
	userInteractor := interactors.NewUserInteractor(userRepo)
	// productInteractor := interactors.NewProductInteractor(productRepo) // Contoh interactor lain

	// API Layer (Routes)
	// Daftarkan semua rute aplikasi dan berikan interactor yang dibutuhkan
	routes.SetupRoutes(app, userInteractor)

	// Mulai server Fiber
	log.Fatal(app.Listen(":" + cfg.Port))
}
