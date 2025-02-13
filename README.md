# Golang API

This is a portfolio project (Working in progress)

## Requirements

**Go 1.23 >**
```
$ Get installer in: https://go.dev/learn/
```

**PostgreSQL**
```
$ Get installers in:
- postgreSQL: https://www.postgresql.org/
- pgAdmin: https://www.pgadmin.org/ (optional)
- dbeaver: https://dbeaver.io/download/ (optional)
```

## Getting Started

All environment variables are in `env.example`, configure them according to your database.

All executables are defined inside `cmd/main.go`.

Run from base folder to envs be read correctly.
```
$ go run cmd/main.go
```

The server should start at: [`http:localhost:8080`](http:localhost:8080)

## Glossary

1. cmd: project start / dependency injection
2. core: business logic / services
3. adapter: responsible of the transformation between a request from the actor to the core
5. repository / storage: As part of the project infrastructure, it is used to store data, create queues or integrate third parties.
## Architecture Goals

1. It follows the principles of Hexagonal Architecture that separates business logic from external dependencies using ports and adapters, making the system technology-independent, testable, flexible and easy to maintain.

## Attention! This project is WIP and might have bugs.
