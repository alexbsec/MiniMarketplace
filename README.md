# Mini Marketplace

This project is a simple backend system that simulates a mini marketplace where users 
can manage products and services. It’s built using Go with PostgreSQL and demonstrates 
clean architecture, logging, and API management.

## Features

- RESTful API for managing products and services.
- PostgreSQL integration for persistent storage.
- Custom logging using Go's slog with colorized output.
- Dockerized setup for easy local development, testing and future deployment.
- Minimalistic architecture for scaling.

## Technologies Used

- Go (Golang): Backend programming language.
- PostgreSQL: Relational database for data persistence.
- Docker & Docker Compose: Containerization of the backend and database service.

## Project Structure


```bash
/MiniMarketplace
    ├── bin
    └── src
        ├── controllers
        ├── core
        ├── db
        │   ├── config
        │   ├── migrations
        │   └── models
        │       ├── tests
        │       └── utils
        ├── logging
        └── scripts
```

## Setup and installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/MiniMarketplace.git
cd MiniMarketplace
```

2. Install dependencies:

### Locally

If running locally, start with:

```bash
go mod tidy
```
Then, you can build using `make`:

```bash
make build
```

This will create a binary file on `./bin`. Running it should start the
application.

### Docker

If runnning with docker, run inside root directory:

```bash
go mod tidy
```

Then, simply run `docker-compose`:

```bash
docker-compose up --build
```

