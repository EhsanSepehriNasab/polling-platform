# Polling Platform

## Overview

The **Polling Platform** is a scalable and interactive system designed to handle high concurrency with support for both mobile and web clients. It allows users to participate in polls, skip polls, and filter them by tags, with a daily voting limit of 100 votes per user.

Certainly! I'll format the concurrency & RPS results into a table and incorporate it into your README under the relevant section. Here's how it would look:


## Concurrency & RPS Results

The following table shows the results from load testing with varying levels of concurrency and the corresponding request rates (RPS), average response times, and 95th percentile response times:

| Stage (Concurrency) | RPS (Requests/sec) | Avg Resp Time (ms) | P95 Resp Time (ms) |
|---------------------|--------------------|--------------------|--------------------|
| 10                  | 9.91               | 6.89               | 24.15              |
| 20                  | 19.88              | 5.07               | 19.52              |
| 50                  | 49.38              | 7.16               | 25.58              |
| 100                 | 98.79              | 9.24               | 44.85              |


Here's the updated `README.md` section with your caching, event streaming, and concurrency strategy integrated, without changing the existing content. I’ve added a new section titled **Advanced Architecture: Caching, Event Streaming & Concurrency** and linked your document within it.

---

## Advanced Architecture: Caching, Event Streaming & Concurrency

To further optimize scalability and responsiveness, the platform integrates advanced caching techniques, event streaming, and concurrency-safe mechanisms.


- **Caching**: Per-user feed and poll results are cached in Redis for faster access. Pattern-based invalidation is used but is being evolved to event-based strategies.
- **Event Streaming**: Kafka (or similar systems like Redis Streams, NATS) can be used to emit and consume key user actions such as `poll_created`, `poll_voted`, and `poll_skipped` for:
  - Real-time cache invalidation
  - Feed pre-warming
  - Analytics logging
  - Real-time user notifications
  - ML-based personalization
- **Concurrency**: Redis Lua scripting and locking (e.g. Redlock) ensure atomic operations for rate-limiting and avoiding cache stampedes.

📄 **Read the full architecture document here**:  
[Polling Platform: Caching, Event Streaming, and Concurrency Strategy](https://docs.google.com/document/d/1zDXjrkMQWaOBtXNQMiIUYAzWCVSAjVMPz3gsK5qNFoo/edit?usp=sharing)

### 2. Poll Service Architecture

![Polling Platform: Caching and Event Streaming](https://github.com/EhsanSepehriNasab/polling-platform/raw/main/Polling_Platform_Caching_and_E.png)

## Key Features

- **Poll Creation**: Users can create polls with multiple-choice options and tags.
- **Poll Feed**: Users can view polls filtered by tags, ensuring no repeat polls.
- **Voting**: Users can vote on a poll, selecting one option per poll.
- **Skipping**: Users can skip polls without voting.
- **Poll Statistics**: Aggregate vote counts for each poll option are available.
- **Rate Limiting**: Users are limited to 100 votes per day.
- **API Documentation**: Interactive API docs are provided via Swagger UI for easy testing and exploration.

## Tech Stack

- **Backend**: Go (Golang) v1.24
- **Database**: PostgreSQL v14 (storing relational data for polls, options, votes, and users)
- **Caching**: Redis v6 (for cache & fast access to frequently queried polls and stats)
- **Metrics & Monitoring**: Prometheus for system observability
- **API Docs**: Swagger UI for API testing and documentation
- **Routing**: Chi Router
- **Containerization**: Docker and Docker Compose for local development environments

---

## Project Structure

```
polling-platform/
├── cmd/
│   └── server/
│       └── main.go          # Main entry point for the server
├── docs/                    # Swagger documentation files
├── internal/
│   ├── db/                  
│   │   └── db.go            # Database connection 
│   │   └── migrate.go       # Database migrations
│   ├── config/              
│   │   └── config.go        # Configuration settings
│   ├── polls/
│   │   ├── handler.go       # HTTP request handlers for polls
│   │   ├── repository.go    # Poll repository (data access layer)
│   │   └── service.go       # Business logic for polls
│   ├── users/
│   │   ├── handler.go       # HTTP request handlers for users
│   │   ├── repository.go    # users repository (data access layer)
│   │   └── service.go       # Business logic for users
│   ├── cache/               
│   │   └── redis.go         # Redis caching for popular data
│   ├── metrics/
│   │   └── prometheus.go    # Prometheus integration for metrics
├── models/
│   └── poll.go              # Poll data model
│   └── user.go              # User data model
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── docker-compose.yml       # Local Docker Compose setup
├── README.md                # Project documentation
├── .env                     # Environment vars
├── concurrency&RPS.txt      # Response times vs. concurrency/RPS Results from k6
├── prometheus.yml           # Prometheus config
├── load_test.js             # Script for load tests
├── Dockerfile               # Dockerfile for polling service
├── load_test_summary.json   # Result of load test


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

2. **Set up environment variables**:
   Create a `.env` file with the following content:
   ```env
    DATABASE_URL=postgres://postgres:postgres@postgres:5432/polling_platform
    PORT=8080
    REDIS_ADDR=redis:6379
    REDIS_PASSWORD=
    REDIS_DB=0
   ```

3. **Start the application with Docker**:
   To start the app and its dependencies (PostgreSQL, Redis, and Prometheus) in Docker containers, run:
   ```bash
   docker-compose up --build -d
   ```

4. **Database Migrations**:  
   Migrations will automatically run on startup to set up the necessary database schema.

5. **API Documentation**:  
   Once the server is running, you can access the interactive Swagger UI to test API endpoints:
   ```
   http://localhost:8080/swagger/index.html
   ```

6. **Prometheus metrics**:  
   Once the server is running, you can access the Prometheus metrics:

   Prometheus panel: 
   ```
   http://localhost:9090/
   ```

   Metrics endpoint: 
   ```
   http://localhost:8080/metrics/
   ```

---

## API Endpoints

### 1. Create Poll

- **Endpoint**: `POST /polls`
- **Request Header**:
  ```
  userId: int
  ```
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
- **Request Header**:
  ```
  userId: int
  ```
- **Query Parameters**:
    - `tag` (optional): Filter by tag
    - `page` (optional): Page number for pagination
    - `limit` (optional): Number of polls per page

- **Response**: Returns a list of polls that the user has not yet voted on or skipped.

### 3. Vote on a Poll

- **Endpoint**: `POST /polls/{pollId}/vote`
- **Request Header**:
  ```
  userId: int
  ```
  
- **Request Body**:
    ```json
    {
        "optionIndex": 0
    }
    ```

- **Response**: `200 OK` or `204 No Content`

### 4. Skip a Poll

- **Endpoint**: `POST /polls/{pollId}/skip`
- **Request Header**:
  ```
  userId: int
  ```

- **Response**: `200 OK` or `204 No Content`

### 5. Poll Statistics

- **Endpoint**: `GET /polls/{pollId}/stats`
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

- **PostgreSQL** stores all relational data related to polls, options, votes, and users.
- **Schema**: 
    - **Polls**: Stores metadata like title and tags.
    - **Options**: Stores the possible choices for each poll.
    - **Votes**: Links individual votes to a user and option.
    - **Users**: Stores user-specific information.
- Migrations are managed automatically.

### Caching

- **Redis** provides in-memory caching for frequently accessed data, such as popular polls and aggregated stats, improving performance.
- Cache invalidation occurs when votes are cast or polls are skipped or poll created.

### API Documentation

- **Swagger UI** offers an interactive interface for exploring and testing the API directly from your browser.

### Metrics & Observability

- **Prometheus** is used for monitoring system performance, such as request counts, latencies, and cache hit/miss ratios. Metrics are exposed at `/metrics` for scraping.

---

## Future Enhancements

1. **Real-Time Updates**: Implement WebSockets or long-polling to deliver live vote updates.
2. **Advanced Poll Types**: Enable multi-stage polls and real-time leaderboards.
3. **Enhanced Caching**: Explore distributed caching for scalability.
4. **Scalability**: Add support for horizontal scaling via technologies like Kubernetes.

---

## Assumptions & Trade-offs

- **Database Choice**: PostgreSQL provides strong consistency and relational query capabilities, but it may face scalability challenges with very high traffic.
- **Caching**: Redis improves read performance but requires careful management of cache invalidation.
- **Rate Limiting**: The current 100 votes per day cap is fixed but could be made more flexible in the future.

---

## Conclusion

The Polling Platform provides a scalable and efficient solution for interactive polling, leveraging PostgreSQL, Redis, and Prometheus to handle high concurrency while maintaining performance and observability.

---

## Future Scaling & Evolution

### Scaling for More Users

- **Horizontal Scaling**: Implement load balancing and microservices to better distribute traffic and isolate services.
- **Caching**: Use CDNs and Redis to cache frequently accessed data.
- **Database Scaling**: Introduce read replicas or consider NoSQL databases like MongoDB for better scalability.

### Supporting Complex Poll Types

- **Multi-Stage Polls**: Allow sequential voting with session tracking.
- **Real-Time Leaderboards**: Use WebSockets to update leaderboards live.

### Future Features

- **Gamification**: Add badges and leaderboards to encourage user engagement.
- **Personalized Polls**: Suggest polls based on user history using machine learning.
- **Poll Results as NFTs**: Reward participants with NFTs for exclusive polls.

---

## Areas for Refactoring & Technical Debt

- **Metrics Collection**: Centralize logic for better maintainability and extend with more granular metrics.
- **Caching**: Improve cache eviction strategies and implement fallback mechanisms.
- **Database Performance**: Optimize slow queries and consider asynchronous job processing.
- **Error Handling**: Implement custom error types and structured logging for improved visibility.
- **Concurrency**: Optimize concurrency handling with worker pools for better throughput.
- **EventSreaming**: Using kafka to write poll reqeust in queue and have better performance in high scaling.