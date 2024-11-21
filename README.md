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

