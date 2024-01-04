# Base stage: Download and cache dependencies
FROM golang:1.20-alpine AS base

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Build stage: Build the application
FROM base AS build

COPY . .
RUN go build -o bot ./cmd/main.go

# Final stage: Create the production image
FROM alpine:latest

WORKDIR /app

# Copy only the built executable from the build stage
COPY --from=build /app/bot /app/bot

# Expose a port
EXPOSE 8000
