# MountAPI 🏔️

**MountAPI** adalah layanan REST API yang menyediakan data komprehensif mengenai gunung-gunung di Indonesia, jalur pendakian, dan pos-pos (checkpoints) yang ada. Dibangun menggunakan Go, Gin Framework, dan PostgreSQL.

## 🚀 Fitur Utama

- **Open Data Gunung Indonesia**: Informasi nama, ketinggian (mdpl), dan lokasi.
- **Manajemen Jalur Pendakian**: Data jalur resmi untuk setiap gunung.
- **Informasi Pos & Jarak**: Detail setiap pos pendakian beserta jarak ke puncak atau antar pos.
- **Admin Management**: Endpoint khusus untuk mengelola data (CRUD) yang diamankan dengan API Key.
- **Bulk Upload**: Fitur upload pos via CSV untuk kemudahan input data.
- **Swagger Documentation**: Dokumentasi API interaktif.

## 🛠️ Tech Stack

- **Global**: [Go (Golang)](https://go.dev/) components
- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **Configuration**: Viper
- **Documentation**: Swagger (Swaggo)

## 📦 Prasyarat

Sebelum menjalankan proyek ini, pastikan Anda telah menginstal:

- [Go](https://go.dev/dl/) (versi 1.25 atau terbaru)
- [PostgreSQL](https://www.postgresql.org/download/)

## 🏃 Cara Menjalankan

1. **Clone repository ini:**

   ```bash
   git clone https://github.com/qsrrvln/mountapi.git
   cd mountapi
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Konfigurasi Environment:**
   Salin file `.env.example` menjadi `.env` dan sesuaikan dengan konfigurasi database Anda.

   ```bash
   cp .env.example .env
   ```

   _Edit file `.env` dan isi `DB_USER`, `DB_PASSWORD`, `DB_NAME`, dll._

4. **Jalankan Aplikasi:**

   ```bash
   go run cmd/api/main.go
   ```

5. **Akses API:**
   Server akan berjalan di `http://localhost:8080` (default).

## 📚 Dokumentasi API

Selengkapnya mengenai endpoint request dan response dapat dilihat melalui Swagger UI:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### Public Endpoints

- `GET /v1/mountains` - Daftar semua gunung
- `GET /v1/mountains/:id` - Detail gunung
- `GET /v1/mountains/:id/routes` - Daftar jalur pendakian gunung tertentu
- `GET /v1/routes/:id` - Detail jalur
- `GET /v1/routes/:id/posts` - Daftar pos pada jalur tertentu

### Admin Endpoints (Require `X-API-Key` Header)

- **Mountains**: `POST`, `PUT`, `PATCH`, `DELETE`
- **Routes**: `POST`, `PUT`, `PATCH`, `DELETE`
- **Posts**: `POST`, `PUT`, `PATCH`, `DELETE`
- **Bulk Upload**: `POST /v1/posts/bulk` (Upload CSV)

## 🗂️ Struktur Project

```
mountapi/
├── cmd/api/        # Entry point aplikasi
├── internal/
│   ├── config/     # Load konfigurasi env
│   ├── database/   # Koneksi database
│   ├── handlers/   # File handler / controller HTTP
│   ├── middleware/ # Middleware (Auth, CORS, etc)
│   └── models/     # Definisi struct database & migrasi
├── docs/           # Generated Swagger docs
└── ...
```

## 📝 Lisensi

[MIT License](LICENSE)
