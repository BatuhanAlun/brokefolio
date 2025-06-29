# Use the official Go image as a base for building
FROM golang:1.23.9 AS builder

# Set the working directory inside the container for the entire module
WORKDIR /app

# Copy go.mod and go.sum from the host's root project directory
# into the container's /app (which is the WORKDIR)
COPY go.mod go.sum ./

# Download Go modules
# This will download all dependencies based on go.mod and go.sum in /app
RUN go mod download

# Copy the entire project source code (including backend, static, templates, etc.)
# from the host's root project directory into /app in the container.
COPY . .

# Build the Go application.
# The path to main.go is now relative to the WORKDIR /app.
# So it becomes ./backend/cmd/brokefolio/main.go
# CGO_ENABLED=0 for static binary, GOOS=linux for container environment.
# -o /app/brokefolio places the final executable directly at /app/brokefolio
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/brokefolio ./backend/cmd/brokefolio/main.go

# --- Production Stage ---
# Use a new, smaller base image for the final runtime
FROM alpine:latest

# Install ca-certificates needed for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory for the final image
WORKDIR /app

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app/brokefolio /app/brokefolio

# Copy static assets and templates
# These are copied from the *host* project root, as the previous COPY . . command
# in the builder stage copied them all into the builder image, but we are copying
# them directly from the host into the new, clean alpine image here.
COPY static/ ./static/
COPY templates/ ./templates/

# If you have an .env file that needs to be copied into the container
# for runtime configuration, uncomment the line below.
# However, it's generally better to use environment variables directly
# in Docker or Docker Compose for production secrets.
# COPY .env ./.env

# Expose the port your application listens on (assuming 8080)
EXPOSE 8080

# Command to run the executable when the container starts
CMD ["/app/brokefolio"]