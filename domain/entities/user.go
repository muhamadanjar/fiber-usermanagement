package entities

import "time"

// User merepresentasikan entitas pengguna dalam domain aplikasi.
// Ini adalah definisi struktur data inti yang digunakan di seluruh aplikasi.
// Tag `gorm` digunakan untuk pemetaan ORM GORM ke kolom database.
// Tag `json` digunakan untuk serialisasi/deserialisasi JSON saat berinteraksi dengan API.
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`             // ID unik pengguna, kunci utama
	Name      string    `json:"name"`                             // Nama pengguna
	Email     string    `json:"email" gorm:"unique"`              // Alamat email pengguna, harus unik
	Password  string    `json:"-"`                                // Kata sandi pengguna, tag `json:"-"` menyembunyikannya dari output JSON
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Waktu pembuatan record, diisi otomatis oleh GORM
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Waktu pembaruan record terakhir, diisi otomatis oleh GORM
}
