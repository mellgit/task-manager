# Task manager

Task manager - this is a service for adding tasks to a queue.


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

What's next?: 

After authentication, you can:
- create a task 
- get
- delete
- change
- get a list of tasks

When a task is created, it sends it to the worker for processing (processing is an imitation, in fact it does nothing except change the status of the task to `done`)



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


**Note:** swagger documentation is available at `http://localhost:3000/swagger/index.html`

## <a name="Struct"></a> Struct project


```
task-manager/
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── up.go                # entry point
├── config.yml
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── auth                  # registration and auth
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── config
│   │   └── config.go
│   ├── db
│   │   └── postgres.go
│   ├── middleware
│   │   └── jwt.go
│   ├── queue                 # producer and consumer kafka
│   │   ├── consumer.go
│   │   └── producer.go
│   ├── task                  # task use case
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   └── worker                 # worker use case
│       ├── model.go
│       ├── repository.go
│       └── service.go
├── main.go
├── migrations
│   └── 20250515135413_init.sql
└── pkg
    └── logger
        └── logger.go
```