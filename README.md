
# HelloServer with New Relic Integration

HelloServer adalah aplikasi server web sederhana berbasis Go yang terintegrasi dengan **New Relic** untuk monitoring **metrics**, **logs**, dan **traces**.

## Fitur

- **Metrics**: Memantau kinerja setiap endpoint termasuk throughput dan waktu respon.
- **Logs**: Semua log aplikasi diteruskan ke **New Relic Logs** dengan konteks transaksi.
- **Traces**: Distributed tracing memungkinkan Anda melacak request melalui aplikasi.

## Prasyarat

1. **Go**  
   Install Go dari [Go Official Website](https://golang.org/doc/install).

2. **New Relic Account**  
   Daftar akun di [New Relic](https://newrelic.com/) dan dapatkan **License Key**.

## Cara Menjalankan Aplikasi

### 1. Clone Repository
Clone repository ke komputer Anda:
```bash
git clone https://github.com/yourusername/helloserver.git
cd helloserver
```

### 2. Pasang Dependency
Pasang semua dependency yang diperlukan menggunakan perintah berikut:
```bash
go get github.com/newrelic/go-agent/v3
go get github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus
```

### 3. Konfigurasi New Relic
Buka file `server.go` dan ganti `YOUR_NEW_RELIC_LICENSE_KEY` dengan **License Key** dari akun New Relic Anda:
```go
newrelic.ConfigLicense("YOUR_NEW_RELIC_LICENSE_KEY")
```

### 4. Jalankan Server
Jalankan aplikasi menggunakan perintah:
```bash
go run server.go
```

### 5. Akses Endpoint
Server akan berjalan di `http://localhost:8080`. Gunakan browser atau `curl` untuk mengakses endpoint berikut:
- **`/`**: Menampilkan pesan sapaan. Contoh:
  ```
  Hello, Gopher!
  ```
- **`/version`**: Menampilkan informasi build aplikasi. Contoh:
  ```
  go 1.20.0
  dep: github.com/sirupsen/logrus v1.9.0
  dep: github.com/newrelic/go-agent/v3 v3.21.0
  ```

## Penjelasan Integrasi dengan New Relic

### 1. **Metrics**
Setiap endpoint dibungkus menggunakan middleware New Relic (`newrelic.WrapHandleFunc`). Ini memungkinkan New Relic mencatat transaksi secara otomatis:
```go
http.HandleFunc(newrelic.WrapHandleFunc(app, "/", greet))
http.HandleFunc(newrelic.WrapHandleFunc(app, "/version", version))
```

### 2. **Logs**
Log aplikasi dikirimkan ke **New Relic Logs** menggunakan `nrlogrus`. Formatter ini memastikan setiap log dikaitkan dengan transaksi yang relevan:
```go
logger.SetFormatter(nrlogrus.NewFormatter(app, &logrus.TextFormatter{}))
```

### 3. **Traces**
Transaksi diinisialisasi dalam setiap handler menggunakan `StartTransaction` dan dikaitkan dengan HTTP request menggunakan `SetWebRequestHTTP`:
```go
txn := app.StartTransaction("Greet")
txn.SetWebRequestHTTP(r)
defer txn.End()
```

## Monitoring di New Relic

### 1. **APM (Application Performance Monitoring)**
Setelah server berjalan, buka **New Relic APM** di dashboard Anda:
1. Pilih aplikasi **HelloServer**.
2. Di tab **Transactions**, Anda dapat melihat transaksi seperti `Greet` dan `Version`.
3. Tab **Metrics** menampilkan throughput dan response time untuk setiap endpoint.

### 2. **Logs**
Masuk ke **New Relic Logs** dan cari log berdasarkan kata kunci, misalnya:
```
"endpoint": "/version"
```

### 3. **Distributed Tracing**
Masuk ke tab **Distributed Traces** di dashboard **APM** untuk melihat detail tracing request.

## Debugging

Jika Anda menghadapi masalah, aktifkan **debug logging**:
```go
newrelic.ConfigDebugLogger(os.Stdout)
```

Setelah diaktifkan, debug log akan muncul di terminal Anda. Contoh:
```
[DEBUG] New Relic Transaction recorded: Name=Greet
```

## FAQ

### Q: Mengapa logs tidak muncul di New Relic Logs?
1. Pastikan Anda sudah mengatur **License Key** dengan benar.
2. Pastikan aplikasi memiliki akses ke domain:
   - `*.newrelic.com`
   - `*.nr-data.net`
3. Cek output debug di terminal.

### Q: Mengapa metrics atau traces tidak muncul di APM?
1. Pastikan setiap transaksi diinisialisasi dengan `StartTransaction` dan dikaitkan dengan HTTP request menggunakan `SetWebRequestHTTP`.
2. Periksa debug log untuk error yang mungkin terjadi.

## Kontribusi

Jika Anda ingin berkontribusi pada proyek ini:
1. Fork repository ini.
2. Buat branch baru untuk perubahan Anda:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Commit dan push perubahan Anda:
   ```bash
   git commit -m "Add my feature"
   git push origin feature/my-feature
   ```
4. Buat Pull Request di GitHub.

## Lisensi

Proyek ini dilisensikan di bawah **MIT License**. Lihat file [LICENSE](LICENSE) untuk detailnya.
