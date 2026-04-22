# Subscription Aggregator Service

REST service designed for aggregating and managing user subscription data.

---

## 🏗 Project Structure

The project follows a clean architecture pattern to ensure separation of concerns and maintainability:

* **`cmd/app`**: Main entry point of the application.
* **`cmd/migrator`**: Utility for executing database migrations.
* **`internal/api`**: Transport layer containing HTTP handlers and routing logic.
* **`internal/service`**: Business logic layer (validation, subscription cost calculation, etc.).
* **`internal/repository`**: Data access layer for interaction with **PostgreSQL**.
* **`internal/model`**: Core domain models and entities.
* **`internal/config`**: Application configuration management.
* **`docs`**: API documentation (Swagger/OpenAPI schemes).
* **`migrations`**: SQL scripts for database schema versioning.

---

## 🚀 Quick Start

Follow these steps to get the service up and running in your local environment:

### 1. Environment Configuration
Create a `.env` file in the root directory. You can use the provided template as a starting point:
```bash
cp .env.example .env
```

### 2. Launch with Docker

Build and start all services (App, Database) using Docker Compose:

```bash
docker-compose up --build
```