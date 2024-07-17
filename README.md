# Simple Go API with CICD 
## Prerequisite
- Github 
- Github-Action
- VPS

## How CICD Works 
0.  Setup Server Agar bisa mendeploy Binary Golang
1. Checkout dan init env golang
2. Menginstall Dependensi Golang 
3. Menjalankan Linter untuk mengecek apakah ada kesalahan dalam standar penulisan code
4. Menjalankan Testing untuk mengecek fungsi yang telah dibuat sudah sesuai dengan ekspektasi 
5. Build execute binary golang
6. Copy binary ke server menggunakan scp 
7. Restart service binary golang dengan mengeksekusi file bash yang telah dibuat di dalam server
8. Untuk mengakses service nya ada pada alamat berikut : http://103.174.114.151:8080/items

## How Its Can Be Restart
- Pada Server, service golang dijalankan pada background proses menggunakan systemd
- Ketika terjadi scp file binary yang lama akan secara otomatis ditimpa
- file bash hanya menjalankan perintah systemctl untuk merestart service yang berjalan dengan systemd

