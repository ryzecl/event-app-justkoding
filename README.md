# 🎟️ Event App - Backend RESTful API

RESTful API backend untuk platform manajemen dan pemesanan tiket acara (*Event Management & Booking System*) yang dibangun menggunakan **Golang**, **Gin Gonic**, **GORM**, dan **PostgreSQL**.

---

## 🚀 Fitur Utama

- 🔐 **Autentikasi & Otorisasi:**
  - Registrasi user dengan validasi email unik.
  - Login dengan hashing password aman (**bcrypt**) dan **JWT (JSON Web Token)**.
  - Endpoint profil user yang terautentikasi (`/auth/me`).
- 🎪 **Manajemen Event:**
  - CRUD (Create, Read, Update, Delete) Event.
  - Upload & validasi gambar ke cloud storage (**ImageKit.io**).
  - Pencarian fleksibel (*case-insensitive search* menggunakan `ILIKE` pada judul & deskripsi).
  - Pembagian halaman (*Pagination*) dinamis (`page` & `limit`) disertai respon `meta`.
  - Hak akses: hanya pembuat event yang dapat mengedit atau menghapus event miliknya.
- 🎫 **Sistem Booking Tiket:**
  - Pemesanan tiket event dengan pembuatan kode booking otomatis (`BK-YYYYMMDD...`).
  - Pencegahan pemesanan ganda (*prevent duplicate booking* untuk event yang sama).
  - Melihat riwayat tiket yang sudah dibooking (*Nested Preload* data event & pembuatnya).
  - Pembatalan/penghapusan booking tiket dengan verifikasi kepemilikan.
- 🌐 **CORS Ready:**
  - Terintegrasi dengan middleware `gin-contrib/cors` agar API dapat diakses dengan lancar dari aplikasi frontend (React, Next.js, Vue, mobile app, dll).

---

## 🛠️ Tech Stack

- **Bahasa:** [Golang](https://go.dev/) (Go 1.22+)
- **Framework Web:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **CORS Middleware:** [gin-contrib/cors](https://github.com/gin-contrib/cors)
- **Database & ORM:** [PostgreSQL](https://www.postgresql.org/) & [GORM](https://gorm.io/)
- **Autentikasi:** [golang-jwt/jwt](https://github.com/golang-jwt/jwt) & [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **File & Media Storage:** [ImageKit Go SDK](https://github.com/imagekit-developer/imagekit-go)
- **Environment Management:** [godotenv](https://github.com/joho/godotenv)

---

## 📁 Struktur Folder

```text
server/
├── config/             # Konfigurasi database & service eksternal
│   ├── db.go           # Koneksi PostgreSQL & AutoMigrate GORM
│   └── imagekit.go     # Inisialisasi SDK ImageKit
├── controllers/        # Logic handler untuk setiap endpoint API
│   ├── booking_controller.go
│   ├── event_controller.go
│   └── user_controller.go
├── middlewares/        # HTTP Middleware (Auth JWT)
│   └── auth_middleware.go
├── models/             # Schema tabel & entity database
│   ├── booking.go
│   ├── event.go
│   └── user.go
├── .env                # Variabel konfigurasi lingkungan (rahasia)
├── .env.example        # Template variabel lingkungan
├── go.mod              # Definisi module & dependensi
└── main.go             # Entry point aplikasi & pendaftaran router
```

---

## ⚙️ Panduan Instalasi & Menjalankan Aplikasi

### 1. Prasyarat

Pastikan kamu telah menginstal:

- [Go](https://go.dev/dl/) (versi 1.22 ke atas)
- [PostgreSQL](https://www.postgresql.org/) (atau cloud database seperti Neon.tech, Supabase)
- Akun [ImageKit.io](https://imagekit.io/) (gratis)

### 2. Masuk ke Folder Project

```bash
cd server
```

### 3. Setup Environment Variable

Duplikasi file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Sesuaikan nilai variabel di dalam file `.env`:

```env
DATABASE_URI="postgresql://username:password@localhost:5432/event_db?sslmode=disable"
JWT_SECRET="kunci_rahasia_jwt_kamu"

# ImageKit Configuration
IMAGEKIT_PUBLIC_KEY="public_xxxx"
IMAGEKIT_PRIVATE_KEY="private_xxxx"
IMAGEKIT_URL_ENDPOINT="https://ik.imagekit.io/your_id/"
```

### 4. Install Dependensi

```bash
go mod download
```

### 5. Jalankan Server

```bash
go run .
# atau
go run main.go
```

Server akan berjalan di `http://localhost:8080`. Database akan otomatis ter-migrasi oleh GORM (`AutoMigrate`).

---

## 📖 Dokumentasi Endpoint API

Semua endpoint diawali dengan prefix: `/api`

### 1. Autentikasi (`/auth`)

| Method   | Endpoint           | Auth | Deskripsi                                      | Body / Parameter                         |
| :------- | :----------------- | :--: | :--------------------------------------------- | :--------------------------------------- |
| `POST` | `/auth/register` |  ❌  | Pendaftaran akun user baru                     | JSON:`{ "name", "email", "password" }` |
| `POST` | `/auth/login`    |  ❌  | Login dan mendapatkan token JWT                | JSON:`{ "email", "password" }`         |
| `GET`  | `/auth/me`       |  ✅  | Mendapatkan data profil user yang sedang login | -                                        |

---

### 2. Events (`/events`)

| Method     | Endpoint         | Auth | Deskripsi                                          | Parameter / Body                                                                           |
| :--------- | :--------------- | :--: | :------------------------------------------------- | :----------------------------------------------------------------------------------------- |
| `GET`    | `/events`      |  ❌  | Mengambil semua event (search & pagination)        | Query:`?search=&page=1&limit=6`                                                          |
| `GET`    | `/events/:id`  |  ❌  | Mengambil detail event beserta daftar peserta      | Param:`id` (Event ID)                                                                    |
| `GET`    | `/events/user` |  ✅  | Mengambil daftar event yang dibuat oleh user login | -                                                                                          |
| `POST`   | `/events`      |  ✅  | Membuat event baru                                 | Multipart Form:`name`, `description`, `location`, `datetime`, `image`            |
| `PUT`    | `/events/:id`  |  ✅  | Mengedit event yang dimiliki                       | Multipart Form:`name`, `description`, `location`, `datetime`, `image` (opsional) |
| `DELETE` | `/events/:id`  |  ✅  | Menghapus event yang dimiliki                      | Param:`id` (Event ID)                                                                    |

---

### 3. Booking (`/booking`)

| Method     | Endpoint          | Auth | Deskripsi                                        | Parameter / Body                              |
| :--------- | :---------------- | :--: | :----------------------------------------------- | :-------------------------------------------- |
| `POST`   | `/booking`      |  ✅  | Memesan tiket untuk suatu event                  | JSON:`{ "phone": "0812...", "eventId": 1 }` |
| `GET`    | `/booking/user` |  ✅  | Mengambil seluruh tiket yang dibooking oleh user | -                                             |
| `DELETE` | `/booking/:id`  |  ✅  | Membatalkan / menghapus pemesanan tiket          | Param:`id` (Booking ID)                     |

---

## 🔒 Catatan Keamanan

- Endpoint dengan label **Auth (✅)** mewajibkan header HTTP:
  ```http
  Authorization: Bearer <token_jwt_kamu>
  ```
- Password user di-hash menggunakan algoritma **bcrypt** sebelum disimpan ke database.
- Sensitif field seperti `password` selalu dikecualikan saat mengembalikan relasi data user (`Preload User`).
