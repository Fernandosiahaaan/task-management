# Collaborative Task Management (Backend) 💻

## 🖊 Overview

Collaborative Task Management System is a platform designed to help users efficiently manage their tasks, whether for personal use or in a team context. This system allows users to create, track, update, and delete tasks in real-time, as well as provides collaboration capabilities among team members, including real-time notifications, task reminders, and activity logging.
The development utilizes several tech stacks, such as:

- Golang + Framework gorilla mux
- Microservice (user service, task service, log service, notification service).
- Clean architecture
- JWT Token Auth
- Redis -> for caching data
- RabbitMq -> message queue for async communication microservice (in progress)
- gRPC (Google Remote Procedure Call) -> for sync communication microservice
- PostgreeDB (SQL DB) -> for task, and user microservice DB
- MongoDB (No SQL DB) -> for logging microservice DB (in progress)
- Docker -> for setup env & bundling all stack tech

## 🖊 Documentation

- Documentation : [Tech Documentaion](https://maroon-crabapple-bb5.notion.site/Collaborative-Task-Management-Backend-1107b515908e80a997c3ee75907ffb2b?pvs=4)

## 🖊 PreRequire

- install [vscode](https://code.visualstudio.com/download)
- install [docker dekstop](https://www.docker.com/products/docker-desktop/) for your OS
- Download [golang binary](https://go.dev/doc/install)

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

## Build Proto for GRPC

```
cd /task-service/internal/gRPC
protoc --go_out=./user --go_opt=paths=source_relative ./proto/*.proto
protoc --go-grpc_out=./user --go-grpc_opt=paths=source_relative ./proto/*.proto

cp -r /user ../../../user-service/internal/gRRPC/
```
