package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware adalah contoh middleware untuk autentikasi.
// Dalam aplikasi produksi nyata, Anda akan mengimplementasikan logika yang lebih kompleks
// seperti memverifikasi JSON Web Token (JWT), memeriksa sesi, atau mengintegrasikan dengan
// penyedia identitas eksternal.
func AuthMiddleware(c *fiber.Ctx) error {
	// Contoh sederhana: Periksa header "Authorization"
	// Anda harus mengganti logika ini dengan metode autentikasi yang aman dan sesuai.
	authHeader := c.Get("Authorization")
	if authHeader == "" || authHeader != "Bearer my-secret-token" { // Ganti "my-secret-token" dengan token nyata dari database/konfigurasi
		log.Println("Akses tidak sah: Token tidak ada atau tidak valid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tidak sah"})
	}

	// Jika autentikasi berhasil, lanjutkan ke handler berikutnya dalam rantai middleware/rute.
	return c.Next()
}
