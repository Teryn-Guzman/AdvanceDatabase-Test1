# Restaurant Reservation Management System API

A RESTful API for managing restaurant reservations, tables, shifts, waitstaff, and customer information. Built with Go using the standard library for HTTP routing.

## Features

- **Customer Management**: Create, read, update, and delete customer records
- **Table Management**: Manage restaurant tables with capacity information
- **Reservation System**: Handle reservation bookings with status tracking (confirmed, cancelled, no-show, completed)
- **Time Slots**: Define available booking time slots
- **Waitstaff Management**: Manage staff assignments
- **Shift Management**: Handle staff work shifts
- **Special Requests**: Store customer special requests (dietary restrictions, celebrations, accessibility needs)
- **Table Assignments**: Link reservations to tables
- **Shift Assignments**: Assign waitstaff to shifts and tables
- **Rate Limiting**: Built-in request rate limiting for API protection
- **CORS Support**: Configurable Cross-Origin Resource Sharing

## Project Structure

```
.
├── cmd/
│   ├── api/                    # API server application
│   │   ├── main.go            # Application entry point
│   │   ├── routes.go          # Route definitions
│   │   ├── middleware.go      # HTTP middleware (rate limiting, CORS)
│   │   ├── server.go          # HTTP server configuration
│   │   ├── customers.go       # Customer endpoints
│   │   ├── reservations.go    # Reservation endpoints
│   │   ├── tables.go          # Table endpoints
│   │   ├── shifts.go          # Shift endpoints
│   │   ├── time_slots.go      # Time slot endpoints
│   │   ├── waitstaff.go       # Waitstaff endpoints
│   │   ├── special_requests.go # Special request endpoints
│   │   ├── reservation_table_assignments.go
│   │   ├── shift_table_assignments.go
│   │   ├── healthcheck.go     # Health check endpoint
│   │   ├── helpers.go         # Helper functions
│   │   └── errors.go          # Error handling
│   └── examples/              # Example applications
│       └── cors/              # CORS configuration examples
├── internal/
│   ├── data/                  # Data access layer
│   │   ├── customers.go
│   │   ├── reservations.go
│   │   ├── tables.go
│   │   ├── shifts.go
│   │   ├── time_slots.go
│   │   ├── waitstaff.go
│   │   ├── special_requests.go
│   │   ├── reservation_table_assignments.go
│   │   ├── shift_table_assignments.go
│   │   ├── filters.go         # Query filtering utilities
│   │   └── errors.go
│   └── validator/             # Input validation
│       └── validator.go
├── migrations/                # Database migrations
│   ├── 000001_create_customers_table.up.sql
│   ├── 000002_create_tables_table.up.sql
│   ├── 000003_create_shifts_table.up.sql
│   ├── 000004_create_time_slots_table.up.sql
│   ├── 000005_create_reservations_table.up.sql
│   ├── 000006_create_reservation_table_assignments_table.up.sql
│   ├── 000007_create_special_requests_table.up.sql
│   ├── 000008_create_waitstaff_table.up.sql
│   └── 000009_create_shift_table_assignments_table.up.sql
├── api/                       # Database schema/API exports
├── Makefile                   # Build and run commands
├── go.mod                     # Go module definition
└── go.sum                     # Dependency checksums
```

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd Advance-Database
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (see Configuration section)

4. Run database migrations:
   ```bash
   make db/migrations/up
   ```

5. Start the server:
   ```bash
   make run/api
   ```

## Configuration

The application is configured via environment variables. Create a `.envrc` file or export the following variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `4000` |
| `ENVIRONMENT` | Runtime environment (`development`, `staging`, `production`) | `development` |
| `RESTAURANT_DB_DSN` | PostgreSQL connection string | `postgres://restaurant:restaurant@localhost/restaurant_management_db` |
| `RATE_LIMITER_RPS` | Requests per second limit | `2` |
| `RATE_LIMITER_BURST` | Rate limiter burst size | `5` |
| `RATE_LIMITER_ENABLED` | Enable/disable rate limiter | `true` |
| `CORS_TRUSTED_ORIGINS` | Space-separated list of trusted CORS origins | empty |

### Example .envrc

```bash
export PORT=4000
export ENVIRONMENT=development
export RESTAURANT_DB_DSN=postgres://restaurant:restaurant@localhost/restaurant_management_db
export RATE_LIMITER_RPS=10
export RATE_LIMITER_BURST=20
export RATE_LIMITER_ENABLED=true
export CORS_TRUSTED_ORIGINS="http://localhost:3000 http://localhost:8080"
```

## Running the Application

### Development

```bash
make run/api
```

### Using Go Run Directly

```bash
go run ./cmd/api \
    -port=4000 \
    -env=development \
    -db-dsn="postgres://restaurant:restaurant@localhost/restaurant_management_db" \
    -limiter-rps=10 \
    -limiter-burst=20 \
    -limiter-enabled=true
```

## API Endpoints

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/health` | Check API health status |

### Observability
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/observability/customers/metrics` | Customer-related metrics (expvar) |

### Customers
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/customers` | List all customers |
| GET | `/v1/customers/:id` | Get customer by ID |
| POST | `/v1/customers` | Create new customer |
| PATCH | `/v1/customers/:id` | Update customer |
| DELETE | `/v1/customers/:id` | Delete customer |

### Tables
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/tables` | List all tables |
| GET | `/v1/tables/:id` | Get table by ID |
| POST | `/v1/tables` | Create new table |
| PATCH | `/v1/tables/:id` | Update table |
| DELETE | `/v1/tables/:id` | Delete table |

### Shifts
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/shifts` | List all shifts |
| GET | `/v1/shifts/:id` | Get shift by ID |
| POST | `/v1/shifts` | Create new shift |
| PATCH | `/v1/shifts/:id` | Update shift |
| DELETE | `/v1/shifts/:id` | Delete shift |

### Time Slots
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/time-slots` | List all time slots |
| GET | `/v1/time-slots/:id` | Get time slot by ID |
| POST | `/v1/time-slots` | Create new time slot |
| PATCH | `/v1/time-slots/:id` | Update time slot |
| DELETE | `/v1/time-slots/:id` | Delete time slot |

### Reservations
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/reservations` | List all reservations |
| GET | `/v1/reservations/:id` | Get reservation by ID |
| POST | `/v1/reservations` | Create new reservation |
| PATCH | `/v1/reservations/:id` | Update reservation |
| DELETE | `/v1/reservations/:id` | Delete reservation |

### Waitstaff
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/waitstaff` | List all waitstaff |
| GET | `/v1/waitstaff/:id` | Get waitstaff by ID |
| POST | `/v1/waitstaff` | Create new waitstaff |
| PATCH | `/v1/waitstaff/:id` | Update waitstaff |
| DELETE | `/v1/waitstaff/:id` | Delete waitstaff |

### Special Requests
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/special-requests` | List all special requests |
| GET | `/v1/special-requests/:id` | Get special request by ID |
| POST | `/v1/special-requests` | Create new special request |
| PATCH | `/v1/special-requests/:id` | Update special request |
| DELETE | `/v1/special-requests/:id` | Delete special request |

### Reservation Table Assignments
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/reservation-table-assignments` | Get assignment by reservation and table |
| GET | `/v1/reservation-table-assignments/list` | List all assignments |
| POST | `/v1/reservation-table-assignments` | Create new assignment |
| DELETE | `/v1/reservation-table-assignments` | Delete assignment by reservation and table |

### Shift Table Assignments
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/shift-table-assignments` | Get assignment by shift and table |
| GET | `/v1/shift-table-assignments/list` | List all assignments |
| POST | `/v1/shift-table-assignments` | Create new assignment |
| PATCH | `/v1/shift-table-assignments` | Update assignment |
| DELETE | `/v1/shift-table-assignments` | Delete assignment |

## Database Schema

### Tables Overview

| Table | Description |
|-------|-------------|
| `customers` | Customer information (name, email, phone) |
| `tables` | Restaurant tables (capacity, location, availability) |
| `shifts` | Work shifts (date, start/end times) |
| `time_slots` | Available booking time slots |
| `reservations` | Customer reservations with status tracking |
| `reservation_table_assignments` | Links reservations to tables |
| `special_requests` | Customer special requests |
| `waitstaff` | Waitstaff information |
| `shift_table_assignments` | Links shifts to waitstaff and tables |

### Reservation Status Values

- `confirmed` - Reservation is confirmed
- `cancelled` - Reservation was cancelled
- `no_show` - Customer did not show up
- `completed` - Reservation has been fulfilled

## Available Make Commands

```bash
make run/api              # Run the API server
make db/psql              # Connect to database with psql
make db/migrations/up     # Run all database migrations
make db/migrations/new name=<name>  # Create new migration files
```

## Database Connection

Connect to PostgreSQL directly:
```bash
make db/psql
```