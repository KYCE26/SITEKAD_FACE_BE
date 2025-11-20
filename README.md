SITEKAD atau Sistem Informasi Tenaga Ahli Daya adalah backend service berbasis Golang dengan bantuan Gin Framework yang digunakan untuk menangani seluruh proses bisnis pada aplikasi SITEKAD. Aplikasi ini dirancang untuk mempermudah proses absensi, patroli, dan pengajuan administratif seperti lembur dan cuti, dengan sistem validasi yang kuat dan alur persetujuan yang jelas.

Aplikasi ini difokuskan untuk absensi karyawan, khususnya dalam meningkatkan akurasi, keamanan, serta efisiensi pencatatan aktivitas karyawan maupun petugas operasional.

🎯 Fitur Utama

✅ 1. Absensi Berbasis QR Code

Pengguna melakukan absensi dengan memindai QR Code yang telah ditentukan lokasi/areanya.

Sistem melakukan validasi lokasi, memastikan pengguna berada dalam radius lokasi absensi yang diizinkan.

Identifikasi pengguna (user identity validation) memastikan hanya pengguna yang terdaftar yang dapat melakukan absensi.

Pencatatan waktu masuk dan keluar, lengkap dengan log lokasi dan id android.

✅ 2. Scan Patroli

Petugas melakukan scan titik patroli (umumnya melalui QR Code).

Sistem mencatat waktu, lokasi, dan identitas petugas.

Memastikan patroli dilakukan sesuai rute yang telah ditentukan.

✅ 3. Pengajuan Lembur + SPL (Surat Perintah Lembur)

Pengguna wajib mengunggah SPL (Surat Perintah Lembur) dalam format gambar atau PDF sebagai syarat pengajuan lembur.

SPL digunakan sebagai bukti dan dokumen formal untuk permohonan lembur.

Admin melakukan approval/penolakan pengajuan lembur melalui dashboard.

Riwayat lembur tersimpan untuk audit dan pelacakan.

✅ 4. Pengajuan Cuti

Pengguna dapat mengajukan cuti (izin, sakit, tahunan, dan lainnya).

Mendukung upload gambar atau lampiran sebagai keterangan tambahan (misalnya surat dokter).

Pengajuan cuti akan melalui proses approval admin sebelum disahkan.

Sistem mencatat kuota cuti dan histori cuti pengguna.

✅ 5. Autentikasi & Otorisasi

Login dengan validasi kredensial.

Menggunakan JWT untuk autentikasi sesi.

Mendukung role admin & user (jika tersedia).

🧰 Teknologi yang Digunakan

Go (Golang)

Gin Web Framework

JWT Authentication

File upload handler (image/pdf)

Modular folder structure (controllers, models, middlewares, helper, config)

🚀 Instalasi & Menjalankan Proyek

1. Clone repository
    ```git clone https://github.com/thxinsomnia/SITEKAD_BE.git```
    ```cd SITEKAD_BE```

2. Install dependencies
   ```Go Mod Tidy```

3. Atur .Env jika Kamu Menggunakan .Env

4. Jalankan Aplikasi
   ```go run main.go```

5. Aplikasi akan berjalan pda port yang ditentukan, Default Port Aplikasi adalah :8080

