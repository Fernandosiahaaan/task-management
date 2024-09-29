# Collaborative Task Management (Backend) 💻

## 🖊 Overview

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

## Documentation

- file tech doc confluence = [Tech Documentaion](https://www.notion.so/Collaborative-Task-Management-Backend-1107b515908e80a997c3ee75907ffb2b?pvs=4)

## 🖊 PreRequire

- install vscode
- install docker dekstop for your OS

## 🖊 Start/ Run project

Runing docker compose

```
cd /task-management
docker-compose up -d
```

After success, run migration

```
cd /migration
go run .
```

After success, run all microservice.

```
# User microservice
cd /user-service
go run .

# task microservice
cd /task-service
go run .

```
