package container

import (
	"fiber-usermanagement/internal/api/handlers"
	"fiber-usermanagement/internal/api/routes"
	"fiber-usermanagement/internal/config"
	"fiber-usermanagement/internal/domain/entities"
	"fiber-usermanagement/internal/infrastructure/database"
	"fiber-usermanagement/internal/infrastructure/persistence"
	"fiber-usermanagement/internal/usecase/interactors"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	DB  *gorm.DB
	App *fiber.App
	Log *zap.Logger
}

// NewContainer menginisialisasi semua dependensi aplikasi
func NewContainer(cfg *config.Config, app *fiber.App, logContainer *zap.Logger) (*Container, error) {
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

	userHandler := handlers.NewUserHandler(userInteractor)

	routeConfig := routes.RouteConfig{
		App:         app,
		UserHandler: userHandler,
	}
	routeConfig.Setup()

	return &Container{
		DB:  db,
		App: app,
		Log: logContainer,
	}, nil
}
