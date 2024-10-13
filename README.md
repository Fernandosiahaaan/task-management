# 📚 Collaborative Task Management (Backend)

## 🖊 Overview

Collaborative Task Management System is a platform designed to help users efficiently manage their tasks, whether for personal use or in a team context. This system allows users to create, track, update, and delete tasks in real-time, as well as provides collaboration capabilities among team members, including real-time notifications, task reminders, and activity logging.
The development utilizes several tech stacks, such as:

- Golang + Framework gorilla mux
- Microservice (user service, task service, log service, notification service).
- Clean architecture
- JWT Token Auth
- Redis caching data -> for caching data
- RabbitMq message queue -> message queue for async communication microservice
- gRPC synchronuis message (Google Remote Procedure Call) -> for sync communication microservice
- PostgreeDB Database (SQL DB) -> for task, and user microservice DB
- MongoDB Database (No SQL DB) -> for logging microservice DB (in progress)
- Docker containers-> for setup env & bundling all stack tech
- DataDog Monitoring tool

## 🖊 Documentation

- Documentation : [Tech Documentation](https://maroon-crabapple-bb5.notion.site/Collaborative-Task-Management-Backend-1107b515908e80a997c3ee75907ffb2b?pvs=4)

## 🖊 Monitoring

- Monitoring tool : [Monitoring Datadog](https://p.us5.datadoghq.com/sb/11f6ed12-8270-11ef-aeeb-36eb61e68aeb-e017eb963d57a887f17717b9c5f3b7e8)

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

### User Proto

```
cd /task-service/infrastructure/gRPC
protoc --go_out=./user --go_opt=paths=source_relative ./proto/user.proto
protoc --go-grpc_out=./user --go-grpc_opt=paths=source_relative ./proto/user.proto

cp -r /user ../../../user-service/infrastructure/gRRPC/
```

### Logging Proto

```
cd /task-service/infrastructure/gRPC
protoc --go_out=./logging --go_opt=paths=source_relative ./proto/log.proto
protoc --go-grpc_out=./logging --go-grpc_opt=paths=source_relative ./proto/log.proto

cp -r /user ../../../user-service/infrastructure/gRRPC/
```
