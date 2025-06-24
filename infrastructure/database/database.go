package database

import (
	"log"

	"gorm.io/driver/mysql" // Contoh: import driver MySQL. Ganti jika Anda menggunakan database lain.
	"gorm.io/gorm"
)

// NewDBConnection membuat koneksi baru ke database menggunakan DSN (Data Source Name) yang diberikan.
// Fungsi ini mengembalikan instance *gorm.DB dan error jika koneksi gagal.
func NewDBConnection(dsn string) (*gorm.DB, error) {
	// Membuka koneksi GORM dengan dialector MySQL.
	// Sesuaikan `mysql.Open(dsn)` dengan dialector database Anda (misal: `postgres.Open(dsn)`, `sqlite.Open("gorm.db")`).
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err // Kembalikan error jika koneksi gagal
	}

	log.Println("Berhasil terhubung ke database!")
	return db, nil // Kembalikan instance GORM DB
}

// Migrate melakukan migrasi skema database secara otomatis.
// Fungsi ini menerima instance *gorm.DB dan variadic interface{} models.
// `models` harus berupa pointer ke struct entitas GORM Anda (misal: &entities.User{}).
func Migrate(db *gorm.DB, models ...interface{}) {
	log.Println("Memulai migrasi database...")
	// `AutoMigrate` akan membuat tabel, kolom, dan indeks yang hilang berdasarkan struct model yang diberikan.
	// Ini tidak akan mengubah tipe kolom yang sudah ada atau menghapus kolom yang tidak terpakai.
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err) // Hentikan aplikasi jika migrasi gagal
	}
	log.Println("Migrasi database berhasil.")
}
