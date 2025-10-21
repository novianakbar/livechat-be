# LiveChat Backend

Aplikasi backend untuk sistem livechat yang mendukung komunikasi real-time antara customer dan agent dengan integrasi AI.

## Konfigurasi AI

Sistem ini mendukung integrasi dengan sistem AI eksternal untuk memproses pesan chat secara otomatis. Berikut adalah langkah-langkah untuk mengaktifkan dan mengkonfigurasi integrasi AI:

### Variable Lingkungan untuk AI

Tambahkan konfigurasi berikut di file `.env`:

```
# Konfigurasi AI
AI_ENABLED=true                              # Aktifkan/nonaktifkan integrasi AI
AI_WEBHOOK_URL=https://your-ai-system.com/webhook  # URL webhook sistem AI eksternal
```

### Alur Integrasi AI

1. **Penerimaan Pesan**: Saat customer mengirim pesan, sistem akan mengirim data ke webhook AI jika `AI_ENABLED=true`.

2. **Format Data yang Dikirim ke AI**:
   ```json
   {
     "chat_id": "uuid-dari-pesan",
     "session_id": "uuid-dari-sesi",
     "user_id": "uuid-dari-pengguna",
     "message": "isi pesan dari customer"
   }
   ```

3. **Autentikasi Sistem AI**:
   - Sistem AI harus login menggunakan endpoint `/api/auth/login` dengan kredensial user role 'ai'
   - Email: `ai@livechat.com`, Password: `password` (ganti dengan password yang aman di production)
   - Sistem akan mendapatkan JWT token untuk digunakan pada request berikutnya

4. **Endpoint Callback untuk Respons AI**:
   - URL: `/api/ai/response`
   - Method: `POST`
   - Headers: `Authorization: Bearer {jwt-token}`
   - Body:
     ```json
     {
       "session_id": "uuid-dari-sesi",
       "message": "respons dari AI"
     }
     ```

### Migrasi Database untuk AI

Jalankan migrasi database untuk menambahkan role 'ai' dan memperbarui constraint pada tabel:

```
$ go run ./cmd/migrate/main.go up
```

Migrasi ini akan melakukan:
1. Menambahkan 'ai' sebagai role yang valid di tabel users
2. Membuat user dengan role 'ai' untuk autentikasi sistem AI
3. Memperbarui constraint pada tabel chat_messages untuk mengizinkan sender_type 'ai'

### Pengujian Integrasi AI

1. Pastikan variabel lingkungan `AI_ENABLED=true` dan `AI_WEBHOOK_URL` sudah dikonfigurasi
2. Mulai aplikasi backend
3. Kirim pesan dari customer melalui widget
4. Periksa log untuk memastikan request dikirim ke webhook AI
5. Uji endpoint callback dengan mengirim respons dari sistem AI

### Keamanan

- Token JWT harus selalu disertakan dalam header `Authorization`
- Ganti password default user AI pada production
- Gunakan HTTPS untuk semua komunikasi dengan sistem AI