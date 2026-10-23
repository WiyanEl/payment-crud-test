# Payment CRUD

Aplikasi pengelolaan data pembayaran berbasis REST API dengan **Golang, Fiber, PostgreSQL, dan React**. Aplikasi menyediakan autentikasi pengguna menggunakan JWT serta fitur CRUD untuk mengelola data pembayaran.

## Fitur Utama

### Autentikasi
- Registrasi pengguna baru.
- Login menggunakan email dan password.
- Password disimpan dalam bentuk hash menggunakan bcrypt.
- Autentikasi API menggunakan JSON Web Token (JWT).
- Endpoint pembayaran hanya dapat diakses oleh pengguna yang sudah login.

### Pengelolaan Pembayaran
- Menampilkan daftar pembayaran dengan pagination.
- Menampilkan detail pembayaran.
- Menambahkan data pembayaran.
- Mengubah data pembayaran.
- Menghapus data pembayaran.
- Validasi data pembayaran.
- Setiap pengguna hanya dapat mengakses data pembayaran miliknya sendiri.

### Antarmuka Pengguna
- Halaman registrasi dan login.
- Dashboard pembayaran.
- Form tambah dan edit pembayaran.
- Konfirmasi penghapusan.
- Tampilan loading, pesan error, dan kondisi data kosong.
- Fitur tampilkan/sembunyikan password.

## Teknologi yang Digunakan

| Bagian | Teknologi |
|---|---|
| Backend | Golang |
| Framework Backend | Fiber |
| Database | PostgreSQL |
| Autentikasi | JWT |
| Password Hashing | bcrypt |
| Frontend | React |
| Build Tool | Vite |
| Komunikasi API | Fetch API |
| Routing Frontend | React Router |

## Persyaratan Sistem

Pastikan perangkat sudah memiliki:

- Go sesuai versi yang ditentukan pada `go.mod`.
- Node.js dan npm.
- PostgreSQL.
- Git.

## Instalasi dan Konfigurasi

### 1. Clone Repository

```bash
git clone <https://github.com/WiyanEl/payment-crud-test>
cd <payment-crud-test>
```

### 2. Buat Database PostgreSQL

Buat database bernama `payment_crud` melalui PostgreSQL atau menggunakan perintah:

```sql
CREATE DATABASE payment_crud;
```

### 3. Jalankan Database Migration

Pastikan file migration tersedia di `migrations/001_init.sql`, kemudian jalankan dari root proyek:

```bash
psql -U postgres -d payment_crud -f migrations/001_init.sql
```

Ganti `postgres` dengan username PostgreSQL yang digunakan jika berbeda.

Migration akan membuat tabel `users` dan `payments`, termasuk relasi pengguna, validasi nilai pembayaran, metode pembayaran, dan status pembayaran.

### 4. Konfigurasi Backend

Buat file `.env` di root proyek dengan konfigurasi berikut:

```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:YOUR_DB_PASSWORD@localhost:5432/payment_crud?sslmode=disable
JWT_SECRET=GANTI_DENGAN_SECRET_YANG_AMAN
```

Sesuaikan username, password, host, dan nama database dengan konfigurasi PostgreSQL milikmu.

**Penting:** Jangan commit file `.env` atau membagikan secret JWT ke repository publik. Pastikan `.env` sudah masuk ke `.gitignore`.

### 5. Jalankan Backend

Dari root proyek, jalankan:

```bash
go mod download
go run .
```

Jika berhasil, backend akan berjalan pada:

```text
http://localhost:8080
```

Endpoint pemeriksaan server:

```text
GET http://localhost:8080/api
```

Respons yang diharapkan:

```json
{
  "success": true,
  "message": "Payment API is running"
}
```

### 6. Konfigurasi Frontend

Buka terminal baru dan masuk ke direktori frontend:

```bash
cd frontend
```

Buat file `.env` di dalam folder `frontend`:

```env
VITE_API_URL=http://localhost:8080/api
```

Instal dependency:

```bash
npm install
```

Jalankan aplikasi:

```bash
npm run dev
```

Buka URL yang ditampilkan Vite, biasanya:

```text
http://localhost:5173
```

Biarkan backend dan frontend tetap berjalan pada terminal masing-masing.

## Cara Menggunakan Aplikasi

### 1. Registrasi

- Buka aplikasi melalui browser.
- Masuk ke halaman registrasi.
- Isi nama, email, dan password.
- Submit formulir registrasi.
- Setelah registrasi berhasil, masuk menggunakan akun tersebut.

### 2. Login

- Masukkan email dan password yang telah didaftarkan.
- Klik tombol login.
- Jika berhasil, pengguna diarahkan ke dashboard pembayaran.

### 3. Melihat Data Pembayaran

Dashboard menampilkan daftar pembayaran milik pengguna yang sedang login, termasuk nomor urut, informasi pelanggan, nominal pembayaran, metode, dan status.

Gunakan pagination untuk berpindah halaman jika data tersedia lebih dari satu halaman.

### 4. Menambahkan Pembayaran

- Klik tombol tambah pembayaran.
- Isi nomor pelanggan.
- Isi nama pelanggan.
- Masukkan nominal pembayaran.
- Pilih metode pembayaran: `transfer`, `ewallet`, atau `va`.
- Pilih status pembayaran jika tersedia pada formulir.
- Simpan data.

Nominal pembayaran harus lebih besar dari nol.

### 5. Mengubah Pembayaran

Klik tombol edit pada data yang ingin diubah, perbarui informasi yang diperlukan, kemudian simpan perubahan.

### 6. Menghapus Pembayaran

Klik tombol hapus pada data yang dipilih dan konfirmasikan tindakan tersebut.

### 7. Logout

Gunakan tombol logout untuk keluar dari aplikasi. Setelah logout, pengguna perlu login kembali untuk mengakses dashboard.

## Dokumentasi REST API

Base URL:

```text
http://localhost:8080/api
```

### Endpoint Autentikasi

| Method | Endpoint | Keterangan | Autentikasi |
|---|---|---|---|
| POST | `/auth/register` | Registrasi pengguna | Tidak |
| POST | `/auth/login` | Login pengguna | Tidak |

### Endpoint Pembayaran

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/payments` | Daftar pembayaran dengan pagination |
| GET | `/payments/:id` | Detail pembayaran |
| POST | `/payments` | Tambah pembayaran |
| PUT | `/payments/:id` | Ubah pembayaran |
| DELETE | `/payments/:id` | Hapus pembayaran |

Seluruh endpoint pembayaran membutuhkan token JWT.

Kirim token melalui header:

```http
Authorization: Bearer <JWT_TOKEN>
```

### Contoh Request Registrasi

`POST /api/auth/register`

```json
{
  "name": "Demo User",
  "email": "demo@example.com",
  "password": "Password123!"
}
```

### Contoh Request Login

`POST /api/auth/login`

```json
{
  "email": "demo@example.com",
  "password": "Password123!"
}
```

Simpan token yang dikembalikan oleh API login dan gunakan token tersebut untuk mengakses endpoint pembayaran.

### Contoh Request Tambah Pembayaran

`POST /api/payments`

```json
{
  "customer_number": "CUST-001",
  "customer_name": "Pelanggan Demo",
  "amount": 150000,
  "method": "transfer",
  "status": "pending"
}
```

Gunakan token JWT pada header `Authorization`.

Nilai metode pembayaran yang tersedia:

- `transfer`
- `ewallet`
- `va`

Nilai status pembayaran yang tersedia:

- `pending`
- `success`
- `failed`

Nama parameter pagination mengikuti implementasi endpoint backend yang digunakan.

## Menjalankan Pengujian

Untuk memeriksa apakah paket backend dapat dikompilasi dan pengujian yang tersedia berjalan, gunakan:

```bash
go test ./...
```

Perintah tersebut hanya menjalankan tes yang sudah tersedia. Pastikan skenario autentikasi, validasi, dan CRUD juga diuji sebelum aplikasi digunakan pada lingkungan produksi.

## Keamanan

- Password di-hash menggunakan bcrypt.
- Endpoint pembayaran dilindungi autentikasi JWT.
- Query pembayaran dibatasi berdasarkan pengguna yang sedang login.
- File `.env` tidak boleh di-commit.
- Untuk produksi, gunakan HTTPS, secret yang kuat, dan konfigurasi database yang aman.

## Pengembang

Dibuat sebagai proyek technical test untuk menunjukkan implementasi REST API, autentikasi, pengelolaan database, serta integrasi frontend dan backend.
