package repositories

import "fiber-usermanagement/internal/domain/entities"

// UserRepository mendefinisikan kontrak (interface) untuk operasi persistensi data User.
// Interface ini menjelaskan *apa* yang bisa dilakukan terhadap data User, tanpa
// peduli *bagaimana* implementasinya (misalnya, menggunakan database, file, atau memori).
// Ini adalah bagian dari layer Domain yang independen dari detail infrastruktur.
type UserRepository interface {
	// Create menambahkan User baru ke penyimpanan. Mengembalikan User yang dibuat atau error.
	Create(user *entities.User) (*entities.User, error)
	// FindByID mencari User berdasarkan ID. Mengembalikan User jika ditemukan atau error jika tidak.
	FindByID(id uint) (*entities.User, error)
	// FindAll mengembalikan semua User yang ada di penyimpanan. Mengembalikan slice User atau error.
	FindAll() ([]entities.User, error)
	// Update memperbarui data User yang sudah ada. Mengembalikan User yang diperbarui atau error.
	Update(user *entities.User) (*entities.User, error)
	// Delete menghapus User berdasarkan ID. Mengembalikan error jika gagal.
	Delete(id uint) error
}
