## Project Overview

> This is a practical and imaginary eCommerce system built with `Go` as the backend. The system is designed using `Monolith Architecture`, `Vertical Slice Architecture`, and `Clean Architecture` principles. The goal is to demonstrate a scalable, maintainable, and testable approach for building modern eCommerce applications.

## Technologies Used

- **[`Go`](https://github.com/golang/go)** - The Go programming language.

- **[`Huma`](https://github.com/danielgtaylor/huma)** - Huma REST/HTTP API Framework for Golang with OpenAPI 3.1.

- **[`Goose`](https://github.com/pressly/goose)** - A database migration tool. Supports SQL migrations and Go functions..

- **[`Go-Chi`](https://github.com/go-chi/chi)** - Lightweight, idiomatic and composable router for building Go HTTP services.

- **[`Go validating`](https://github.com/RussellLuo/validating)** - A Go library for validating structs, maps and slices.

- **[`Stoplight Elements`](https://github.com/stoplightio/elements)** - Build beautiful, interactive API Docs with embeddable React or Web Components, powered by OpenAPI and Markdown.

- **[`Uber Go Zap`](https://github.com/uber-go/zap)** - Blazing fast, structured, leveled logging in Go.

- **[`Uber Go Fx`](https://github.com/uber-go/fx)** - A dependency injection based application framework for Go.

- **[`gocache`](https://github.com/eko/gocache)** - A complete Go cache library that brings you multiple ways of managing your caches.

- **[`go-redis`](https://github.com/redis/go-redis)** - Redis Go client.

- **[`copier`](https://github.com/jinzhu/copier)** - Copier for golang, copy value from struct to struct and more.

- **[`pgx`](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit for Go.

- **[`testify`](https://github.com/stretchr/testify)** - A toolkit with common assertions and mocks that plays nicely with the standard library.

- **[`gofakeit`](https://github.com/brianvoe/gofakeit)** - Random fake data generator written in go.

- **[`viper`](https://github.com/spf13/viper)** - Go configuration with fangs.

- **[`Testcontainers`](https://github.com/testcontainers/testcontainers-go)** - Testcontainers for Go is a Go package that makes it simple to create and clean up container-based dependencies for automated integration/smoke tests.

## Quick Start

### Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.23.4**

Once you have Go 1.23.4 installed, you can set up the project by following these steps:

Clone the repository:

```bash
git clone https://github.com/tguankheng016/golang-ecommerce-monolith.git
```

Run the development server:

```bash
go mod tidy
cd cmd/app && go run main.go
```
