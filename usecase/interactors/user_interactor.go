package interactors

import (
	"errors"

	"fiber-usermanagement/domain/entities"
	"fiber-usermanagement/domain/repositories"

	"gorm.io/gorm"
)

// UserInteractor adalah use case untuk operasi terkait entitas User.
// Ini mengimplementasikan logika bisnis yang berinteraksi dengan UserRepository.
type UserInteractor struct {
	userRepo repositories.UserRepository // Dependensi ke interface UserRepository
}

// NewUserInteractor membuat instance baru dari UserInteractor.
// Menerima implementasi UserRepository untuk dipasangkan.
func NewUserInteractor(ur repositories.UserRepository) *UserInteractor {
	return &UserInteractor{userRepo: ur}
}

// CreateUser adalah use case untuk membuat pengguna baru.
// Ini menangani validasi input dasar dan memanggil repository untuk persistensi.
func (i *UserInteractor) CreateUser(user *entities.User) (*entities.User, error) {
	// Contoh logika bisnis: validasi sederhana
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email dan password tidak boleh kosong")
	}
	// TODO: Dalam aplikasi nyata, Anda harus melakukan hashing password di sini
	// Contoh: user.Password, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Panggil repository untuk menyimpan data
	return i.userRepo.Create(user)
}

// GetUserByID adalah use case untuk mendapatkan pengguna berdasarkan ID.
func (i *UserInteractor) GetUserByID(id uint) (*entities.User, error) {
	// Panggil repository untuk mengambil data
	user, err := i.userRepo.FindByID(id)
	if err != nil {
		// Contoh: Jika error menunjukkan record tidak ditemukan, berikan error yang lebih spesifik
		if errors.Is(err, gorm.ErrRecordNotFound) { // Asumsikan GORM mengembalikan error ini
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}
	return user, nil
}

// GetAllUsers adalah use case untuk mendapatkan semua pengguna.
func (i *UserInteractor) GetAllUsers() ([]entities.User, error) {
	// Panggil repository untuk mengambil semua data
	return i.userRepo.FindAll()
}

// UpdateUser adalah use case untuk memperbarui pengguna.
// Ini mengambil pengguna yang ada, memperbarui bidang yang diizinkan, dan menyimpan perubahan.
func (i *UserInteractor) UpdateUser(id uint, user *entities.User) (*entities.User, error) {
	// Ambil pengguna yang ada terlebih dahulu
	existingUser, err := i.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan untuk diperbarui")
		}
		return nil, err
	}

	// Perbarui hanya field yang diizinkan oleh logika bisnis
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	// TODO: Handle password update secara terpisah dengan hashing dan validasi tambahan

	// Panggil repository untuk menyimpan pembaruan
	return i.userRepo.Update(existingUser)
}

// DeleteUser adalah use case untuk menghapus pengguna.
// Ini dapat mencakup logika bisnis pra-penghapusan, seperti memeriksa dependensi.
func (i *UserInteractor) DeleteUser(id uint) error {
	// Contoh logika bisnis: periksa apakah pengguna memiliki relasi yang tidak boleh dihapus
	// Misalnya, jika pengguna memiliki pesanan aktif, mungkin tidak bisa dihapus.

	// Panggil repository untuk menghapus data
	err := i.userRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pengguna tidak ditemukan untuk dihapus")
		}
		return err
	}
	return nil
}
