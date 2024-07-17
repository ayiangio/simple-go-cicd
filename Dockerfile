# Gunakan image Golang sebagai base image
FROM golang:1.20-alpine

# Buat direktori untuk aplikasi
WORKDIR /app

# Copy semua kode sumber ke WORKDIR
COPY . .

# Download dependencies
RUN go mod download

# Build aplikasi
RUN go build -o /simple-api

# Ekspose port 8080
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["/simple-api"]