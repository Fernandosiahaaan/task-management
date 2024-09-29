# Collaborative Task Management (Backend) 💻

## Overview 🖊

Collaborative Task Management System adalah sebuah platform yang dirancang untuk membantu pengguna mengelola tugas-tugas mereka secara efisien, baik untuk penggunaan pribadi maupun dalam konteks tim. Sistem ini memungkinkan pengguna untuk membuat, melacak, memperbarui, dan menghapus tugas secara real-time, serta memberikan kemampuan kolaborasi antara anggota tim, termasuk notifikasi real-time, pengingat tugas, dan pencatatan aktivitas.
Pengembangan yang dilakukan menggunakan beberapa stack tech, seperti :

- Golang + Framework gorilla mux
- Microservice (user service, task service, log service, notification service).
- Clean architecture
- JWT Token Auth
- Redis -> for caching data
- RabbitMq -> message queue for async communication microservice (in progress)
- gRPC (Google Remote Procedure Call) -> for sync communication microservice
- PostgreeDB (SQL DB) -> for task, and user microservice DB
- MongoDB (No SQL DB) -> for logging microservice DB
- Docker -> for setup env & bundling all stack tech

## PreRequire 🖊

- install vscode
- install docker dekstop for your OS

## Microservice Architecture 🖊

### User Microservicce 📌

User service akan menghandle terkait dari proses login, loguout, update user, dll.

#### DB Schema 🛠

| Column Name | Data Type    | Constraints                   | Description                                 |
| ----------- | ------------ | ----------------------------- | ------------------------------------------- |
| id          | UUID         | PRIMARY KEY, UNIQUE, NOT NULL | Unique identifier for each user             |
| username    | VARCHAR(50)  | UNIQUE, NOT NULL              | Username for user login                     |
| password    | VARCHAR(100) | NOT NULL                      | password for user authentication            |
| email       | VARCHAR(100) | NOT NULL                      | User's email address                        |
| role        | VARCHAR(100) | NOT NULL                      | User's role (e.g., admin, user, superadmin) |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP     | Timestamp when the user was created         |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP     | Timestamp when the user was last updated    |

#### API Endpoint ♻️

| Method | URI              | Description                                               |
| ------ | ---------------- | --------------------------------------------------------- |
| POST   | `/user/register` | Mendaftarkan pengguna baru.                               |
| POST   | `/user/login`    | Mengotentikasi pengguna dan mengembalikan token JWT.      |
| POST   | `/user/logout`   | Logout pengguna dengan menghapus token sesi.              |
| GET    | `/user/aboutme`  | Mendapatkan informasi tentang pengguna yang sedang login. |
| GET    | `/users`         | Mendapatkan daftar semua pengguna .                       |
| PUT    | `/user/update`   | Memperbarui informasi pengguna yang sedang login .        |
