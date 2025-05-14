# Task manager

Task manager - this is a service for adding tasks to a queue.

**Note:** this documentation in development 

## Table of Contents
- [Docker Installation](#Docker)
- [How It Works](#Jobs)
- [Stack](#Stack)
- [Struct project](#Struct)

## <a name="Docker"></a> Docker Installation
### Run
```
make up
```
### Volumes
[Configuration file](./config.yml) `/path/config.yml:/home/app/config.yml:ro`

[Environment file](./.env) `/path/.env:/home/app/.env:ro`

### Compose

The `docker-compose.yml` file contains all the necessary databases

## <a name="Jobs"></a> How It Works

The service implements JWT authentication with access and refresh tokens, with storage of refresh tokens in the database:
- `/register` - registration
- `/login` - we get two access tokens (lives for 5 minutes) and refresh (lives for one week)
- `/refresh` - when the access token expires
- `/logout` - log out (refresh token is deleted from the database)

Additionally:
- if the refresh token expires, you need to get a new one via `/login`
- when using `/refresh` and `/logout`, the request body must contain the `refresh_token` field with the `Bearer` prefix
  (json example: `{"refresh_token": "Bearer <refresh_token>"}`)
- the `/logout` refresh token will be deleted from the database, therefore, use `/login` to receive a new refresh token


**Note:** more details are described in swagger

What's next?: **in development**



## <a name="Stack"></a> Stack

Backend
- **Golang**
- **Fiber**
- **Validator:**
- **goose:** migrations
- **JWT**
- **swagger**

Broker
- **Kafka**

Data Base
- **PostgreSQL**
- **Redis**


**Note:** swagger documentation is available at `http://localhost:3000/swagger/index.html`

## <a name="Struct"></a> Struct project


```

task-manager/
├── cmd/                    # entry point
│   └── main.go
├── internal/
│   ├── task/               # use case
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── queue/              # queue and workers
│   │   ├── producer.go
│   │   └── consumer.go
│   ├── user/               # registration and auth
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   └── repository.go
│   └── middleware/         # JWT, log
├── pkg/
│   └── config/             # env-variable
├── migrations/             # SQL migrations
├── go.mod
└── README.md

```