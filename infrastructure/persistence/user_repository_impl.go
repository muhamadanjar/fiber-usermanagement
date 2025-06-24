package persistence

import (
	"gorm.io/gorm"

	"fiber-usermanagement/domain/entities"
	"fiber-usermanagement/domain/repositories"
)

// UserRepositoryImpl adalah implementasi konkret dari interface repositories.UserRepository.
// Ini menggunakan GORM untuk berinteraksi dengan database.
type UserRepositoryImpl struct {
	db *gorm.DB // Kumpulan koneksi database GORM
}

// NewUserRepository membuat instance baru dari UserRepositoryImpl.
// Ini menerima instance *gorm.DB untuk melakukan operasi database.
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create mengimplementasikan metode Create dari UserRepository.
// Ini membuat record pengguna baru di database.
func (r *UserRepositoryImpl) Create(user *entities.User) (*entities.User, error) {
	result := r.db.Create(user) // GORM akan mengisi ID setelah pembuatan berhasil
	return user, result.Error
}

// FindByID mengimplementasikan metode FindByID dari UserRepository.
// Ini mencari record pengguna berdasarkan ID.
func (r *UserRepositoryImpl) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	result := r.db.First(&user, id) // Mencari record pertama yang cocok dengan ID
	return &user, result.Error
}

// FindAll mengimplementasikan metode FindAll dari UserRepository.
// Ini mengembalikan semua record pengguna dari database.
func (r *UserRepositoryImpl) FindAll() ([]entities.User, error) {
	var users []entities.User
	result := r.db.Find(&users) // Mengambil semua record
	return users, result.Error
}

// Update mengimplementasikan metode Update dari UserRepository.
// Ini memperbarui record pengguna yang sudah ada di database.
func (r *UserRepositoryImpl) Update(user *entities.User) (*entities.User, error) {
	// `Save` akan melakukan operasi update jika record dengan ID tersebut sudah ada,
	// atau insert jika belum ada (upsert). Pastikan `user.ID` diset.
	result := r.db.Save(user)
	return user, result.Error
}

// Delete mengimplementasikan metode Delete dari UserRepository.
// Ini menghapus record pengguna berdasarkan ID.
func (r *UserRepositoryImpl) Delete(id uint) error {
	// Menghapus record User berdasarkan ID. Menggunakan &entities.User{} sebagai model.
	result := r.db.Delete(&entities.User{}, id)
	return result.Error
}
