```
golang-echo-postgresql/
│── main.go                  # Entry point aplikasi
│── go.mod                   # Modul Go untuk dependensi
│── go.sum                   # Checksum dependensi
│── config/                  # Konfigurasi aplikasi
│   ├── config.go            # Konfigurasi untuk koneksi DB dan lainnya
│── db/                      # Folder untuk migrasi database
│   ├── migrations/          # Skrip migrasi SQL untuk pembuatan dan penghapusan tabel
│   │   ├── 001_create_nasabah_table.sql  # Skrip untuk membuat tabel nasabah
│   │   ├── 002_down.sql     # Skrip untuk rollback migrasi
│   ├── db.go                # Koneksi database dan fungsi inisialisasi
│── handlers/                # Handler untuk HTTP request
│   ├── nasabah_handler.go   # Handler untuk operasi CRUD nasabah
│   ├── tabung_handler.go    # Handler untuk operasi CRUD tabung
│── models/                  # Struktur model untuk data
│   ├── nasabah.go           # Definisi model untuk tabel nasabah
│── repositories/            # Repository untuk query database
│   ├── nasabah_repository.go # Repository untuk query data nasabah
│── routes/                  # Rute API
│   ├── routes.go            # Setup dan definisi semua rute
│── utils/                   # Utilitas umum
│   ├── response.go          # Format response standar untuk API
│── .env                     # Environment variables untuk konfigurasi sensitif (DB user, password, dll.)
│── .gitignore               # Mengabaikan file yang tidak perlu di-commit
│── README.md                # Dokumentasi untuk project
│── docker-compose.yml       # File Docker Compose untuk menjalankan DB dan aplikasi
│── Dockerfile               # Dockerfile untuk membangun image aplikasi Golang


    
    ```