# Use the official Go image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp

# Expose a port (if your Go program listens on a specific port)
EXPOSE 8080

# Define the command to run your Go application
CMD ["./myapp"]