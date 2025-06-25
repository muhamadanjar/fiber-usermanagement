package container

import (
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/infrastructure/database"
	"fiber-usermanagement/internal/infrastructure/persistence"
	"fiber-usermanagement/internal/usecase/interactors"
	"log"

	"gorm.io/gorm"
)

// Container struct menyimpan semua dependensi yang diinisialisasi
type Container struct {
	DB             *gorm.DB
	UserInteractor *interactors.UserInteractor
	// Tambahkan interactor lain di sini jika ada
}

// NewContainer menginisialisasi semua dependensi aplikasi
func NewContainer(cfg *config.Config) (*Container, error) {
	// Inisialisasi koneksi database
	db, err := database.NewPostgresDB(cfg.Database.DatabaseURL)
	if err != nil {
		return nil, err
	}
	log.Println("Database connection initialized.")

	// Infrastructure Layer (Persistence)
	userRepo := persistence.NewUserRepository(db)
	log.Println("Repositories initialized.")

	// Usecase Layer (Interactors)
	userInteractor := interactors.NewUserInteractor(userRepo)
	log.Println("Interactors initialized.")

	database.Migrate(db, &entities.User{}, &entities.Role{}, &entities.Permission{})
	log.Println("Migrasi database berhasil.")

	return &Container{
		DB:             db,
		UserInteractor: userInteractor,
	}, nil
}
