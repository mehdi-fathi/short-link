#!/bin/bash

# Define the container name
container_name="go_app"

# Send SIGINT signal to the specific container
echo "Sending SIGINT stop to container $container_name for stop app"

docker compose -f docker/docker-compose.yml --env-file .env.local exec app pkill -SIGINT main

# Wait for a short period to allow containers to handle the SIGINT
echo "Waiting for app to doing all process gracefully..."
sleep 30

# Stop all containers
echo "Stopping containers..."
docker compose -f docker/docker-compose.yml --env-file .env.local stop
