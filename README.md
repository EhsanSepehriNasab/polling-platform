
# Polling Platform

## Overview

The **Polling Platform** is a massively interactive polling system designed to handle high read/write concurrency, with support for mobile and web clients. The platform allows users to vote on polls, skip polls, and filter them by tags. It also enforces a daily vote limit of 100 votes per user.

## Features Implemented

1. **Poll Creation**: Users can create polls with multiple-choice options and tags.
2. **Poll Feed**: Users can view a feed of polls filtered by tags, without seeing the same poll twice.
3. **Vote**: Users can vote on a poll, selecting one option. A vote is recorded for each poll.
4. **Skip**: Users can skip a poll without voting.
5. **Poll Statistics**: Aggregate vote counts per poll option can be retrieved.
6. **Rate Limiting**: Users can vote on up to 100 polls per day.
7. **API Documentation**: Swagger-based documentation for testing the API endpoints.

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL (for relational data management, including polls, options, votes, and users)
- **Caching**: Redis (for in-memory caching of popular polls and stats)
- **Metrics & Monitoring**: Prometheus for observability
- **API Documentation**: Swagger UI for interactive API documentation
- **Routing**: Chi Router
- **ORM**: GORM for PostgreSQL interactions and migrations
- **Testing**: Unit tests and integration tests for business logic and database interactions
- **Containerization**: Docker and Docker Compose for containerized services

---

## Project Structure

```
polling-platform/
├── cmd/
│   └── server/
│       └── main.go          # Main entry point for the server
├── docs/                    # Swagger docs
├── internal/
│   ├── db/                  
│   │   └── db.go            # Database connection & migrations
│   ├── config/                  
│   │   └── config.go        # Setting config
│   ├── polls/
│   │   ├── handler.go       # HTTP request handlers for polls
│   │   ├── repository.go    # Poll repository (data access layer)
│   │   └── service.go       # Business logic for handling polls
│   ├── cache/               
│   └   └── redis.go         # Redis caching for frequently accessed data
│   └── metrics/
│   └  └── prometheus.go     # Prometheus integration for metrics
├── └── models/
│       └── poll.go          # Poll data model
│       
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── docker-compose.yml       # Docker Compose setup for local development
├── README.md                # Project documentation
```

---

## How to Run the Project

### Prerequisites

- Docker and Docker Compose installed
- Go 1.18+ installed

### Step-by-Step Instructions

1. **Clone the repository**:
   ```bash
   git clone https://github.com/EhsanSepehriNasab/polling-platform.git
   cd polling-platform
   ```

2. **Setup environment variables**:
   Create a `.env` file with the following environment variables:
   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@localhost:5432/polling_platform?sslmode=disable
   REDIS_URL=localhost:6379
   ```

3. **Start the application using Docker**:
   Run the following command to start the app along with its dependencies (PostgreSQL, Redis, and Prometheus) in Docker containers:
   ```bash
   docker-compose up --build
   ```

4. **Migrations**:  
   Migrations will be automatically run upon startup to create the necessary database schema.

5. **API Documentation**:  
   Once the server is running, you can access the Swagger UI to test your API endpoints:
   ```
   http://localhost:8080/swagger/index.html
   ```

---

## API Endpoints

### 1. Create Poll

- **Endpoint**: `POST /polls`
- **Request Body**:
    ```json
    {
        "title": "Your favorite programming language?",
        "options": ["Go", "Python", "Rust"],
        "tags": ["programming", "favorites"]
    }
    ```

- **Response**: `201 Created`

### 2. Retrieve Polls for Feed

- **Endpoint**: `GET /polls`
- **Query Parameters**:
    - `tag` (optional): Filter by tag
    - `page` (optional): Page number for pagination
    - `limit` (optional): Number of polls per page
    - `userId` (required): The ID of the user requesting the feed

- **Response**: Returns a list of polls that the user has not yet voted or skipped.

### 3. Vote on a Poll

- **Endpoint**: `POST /polls/{id}/vote`
- **Request Body**:
    ```json
    {
        "userId": 999,
        "optionIndex": 1
    }
    ```

- **Response**: `200 OK` or `204 No Content`

### 4. Skip a Poll

- **Endpoint**: `POST /polls/{id}/skip`
- **Request Body**:
    ```json
    {
        "userId": 999
    }
    ```

- **Response**: `200 OK` or `204 No Content`

### 5. Poll Statistics

- **Endpoint**: `GET /polls/{id}/stats`
- **Response**:
    ```json
    {
        "pollId": 123,
        "votes": [
            { "option": "Go", "count": 10 },
            { "option": "Python", "count": 25 },
            { "option": "Rust", "count": 7 }
        ]
    }
    ```

---

## Technical Details

### Database

- **PostgreSQL** is used to store data regarding polls, options, votes, and user activity.
- **Schema**: 
    - **Polls**: Stores poll metadata (title, tags).
    - **Options**: Stores options for each poll.
    - **Votes**: Stores individual votes, linked to a user and an option.
    - **Users**: Stores user information.
- **Migrations** are handled using `golang-migrate`.

### Caching

- **Redis** is used for caching frequently accessed data like popular polls and aggregated stats, improving response time for these operations.
- Cache invalidation is handled during vote and skip operations.

### API Documentation

- **Swagger** is used for generating interactive API documentation. This allows for testing and exploring the API directly from the browser.

### Metrics & Observability

- **Prometheus** is used for gathering performance metrics (request counts, latencies, cache hits/misses, etc.) to monitor the health and performance of the system.
- Metrics are exposed at the `/metrics` endpoint, where Prometheus can scrape them.

---

## Future Enhancements

1. **Real-Time Updates**: Integrate with WebSockets or long-polling to provide real-time vote updates to users.
2. **Advanced Poll Types**: Allow multi-stage polls (e.g., users vote on multiple questions sequentially) and real-time leaderboards.
3. **Improved Caching**: Use a distributed caching solution to scale with a growing user base.
4. **Scalability**: Add support for horizontal scaling (e.g., using Kubernetes for auto-scaling).

---

## Assumptions & Trade-offs

- **Database Choice**: PostgreSQL was chosen for its strong consistency guarantees and efficient relational queries, though it may not be as scalable as NoSQL solutions for some use cases.
- **Caching**: Redis is used to improve read performance but adds complexity around cache invalidation.
- **Rate Limiting**: A hard limit of 100 votes per day was enforced on users to ensure fair usage. This could be made more flexible in the future (e.g., adding per-user rate limits based on activity).

---

## Running Tests

To run the tests, you can use the `go test` command:

```bash
go test ./...
```

---

## Conclusion

This project implements a highly interactive polling platform that scales efficiently while ensuring data consistency and performance. It uses a combination of PostgreSQL, Redis, and Prometheus to handle high concurrency and maintain robust observability.

---