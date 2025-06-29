# Use the official Go image as a base for building
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's caching.
# This means if only your source code changes, but not dependencies,
# Docker won't re-download modules.
COPY backend/go.mod backend/go.sum ./backend/

# Go into the backend directory within the container
WORKDIR /app/backend

# Download Go modules
RUN go mod download

# Copy the rest of your backend source code
# This copies everything from your local 'backend' folder into '/app/backend' in the container
COPY backend/ .

# Build the Go application
# CGO_ENABLED=0 is important for creating a statically linked binary, which is good for smaller images.
# GOOS=linux ensures it's built for a Linux environment, which is what Docker containers typically are.
# -o /app/brokefolio specifies the output executable name and path (will be in the root of /app)
WORKDIR /app # Change back to /app before building to place the executable directly there
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/brokefolio ./backend/cmd/brokefolio/main.go

# --- This is a multi-stage build. We'll use a new, smaller base image for the final runtime ---
FROM alpine:latest

# Install ca-certificates needed for HTTPS requests from inside the container
RUN apk --no-cache add ca-certificates

# Set the working directory for the final image
WORKDIR /app

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app/brokefolio /app/brokefolio

# Expose the port your application listens on (assuming 8080)
EXPOSE 8080

# Command to run the executable when the container starts
CMD ["/app/brokefolio"]