# Makefile for opggvisualizer with Docker Compose

# Variables
IMAGE_NAME = opggvisualizer
VOLUME_PATH = ./data.db

# Build the Docker images
build:
	docker-compose build

# Initialize the system by fetching initial data
init:
	docker-compose run --rm opggvisualizer champions fetch
	docker-compose run --rm opggvisualizer games fetch

# Start the Docker containers
up:
	docker-compose up -d

# Stop the Docker containers
down:
	docker-compose down

# Clean up Docker resources
clean:
	docker-compose down -v --rmi all

# Help target to display available commands
help:
	@echo "Available commands:"
	@echo "  make build    - Build the Docker images"
	@echo "  make init     - Initialize the system by fetching initial data"
	@echo "  make up       - Start the Docker containers"
	@echo "  make down     - Stop the Docker containers"
	@echo "  make clean    - Clean up Docker resources"
