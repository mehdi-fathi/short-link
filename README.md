# Short-Link

## Project Overview
Short-Link is a URL shortening service created with Golang. It provides a simple and efficient way to shorten long URLs. The project utilizes PostgreSQL for data storage and Redis for caching, offering a robust and scalable solution for URL management.

It's a project to master my skills, in other words an interesting challenge. 
I try to take advantage of substantial concept Golang, engineering and well-structure software engineering.

## Features
- **URL Shortening**: Convert long URLs into shorter, more manageable ones.
- **Analytics**: Track the usage of shortened URLs.
- **Efficient Storage**: Utilize PostgreSQL for storing URL data.
- **Performance**: Leverage Redis for enhanced caching and retrieval speed.

## Tech Features

- Shutdown gracefully GO app with all containers
- Up project with docker
- Take advantage of go routine in saving stat and validate links
- Use redis for saving count visits
- Event-Driven-Design: We used queue with go routines to validate all links right after create a new one. 
- Hexagonal Architecture
- Integration tests

## Installation Instructions

### Prerequisites
- Golang
- PostgreSQL
- Redis
- RabbitMq

### Steps to run

1. **Clone the Repository**:
2. `cp .env.example .env.local`.
3. `make build-docker` for first time.
4. `make up`
5. `make migration_up_v2`

### Step to shut down gracefully

- `make down`
