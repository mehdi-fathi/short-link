

# Use the official Golang image
FROM golang:1.21

# Install netcat
RUN apt-get update && apt-get install -y netcat-traditional


# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY .. .
# Copy go.mod and go.sum files
#COPY ../../go.mod ./

# Download all dependencies. This is done early to take advantage of Docker's caching mechanism.
RUN #go mod download

ENV GO111MODULE=on

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

RUN ["ls"]
RUN ["pwd"]

# Build the Go application
RUN go build -o main ./cmd/main.go

# Command to run the application
CMD ["./main"]

