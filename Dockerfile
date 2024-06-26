

# Use the official Golang image
FROM golang:1.21

# Install netcat
RUN apt-get update && apt-get install -y netcat

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. This is done early to take advantage of Docker's caching mechanism.
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the wait-for-it.sh script
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# Build the Go application
RUN go build -o main ./cmd/*.go


# Command to run the application
CMD ["/wait-for-it.sh", "rabbitmq", "5672", "/wait-for-it.sh", "db", "5432", "/bin/sh", "-c", "migrate -path ./migrations -database postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable up && ./main"]

