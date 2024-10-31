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
- Up project with docker for local mode
- Take advantage of go routine in saving stat and validate links
- Use redis for saving count visits
- Event-Driven-Design: We used queue with go routines to validate all links right after create a new one. 
- Hexagonal Architecture
- Integration tests
- Implement kubernetes for production(helm-nginx)

## Installation Instructions

### Prerequisites
- Golang
- PostgreSQL
- Redis
- RabbitMq
- Grafana

### Steps to run

1. **Clone the Repository**:
2. Local(docker-compose)
   1. `cp .env.example .env.local`.
   2. `make build-docker` for first time.
   3. `make up`
   4. `make migration_up_v2`
   5. Open `http://localhost:8080/index`
3. Kubernetes
   1. `cd deploy/kuber`
   2. `make install_all_cluster`
   3. `make upgrade_app` - After any update in order to create a new version and upgrade.
   4. Run `kubectl get pods` to make sure all pods are running.
   5. Set `127.0.0.1 shortlink.com` in etc/hosts.
   6. Open `https://shortlink.com`

`http://localhost:3000/` is Grafana Dashboard. 

### Step to shut down gracefully in local side (docker)

- `make down`
