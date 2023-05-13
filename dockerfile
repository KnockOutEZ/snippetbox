# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the application
RUN go build ./cmd/web

# Expose port 8080 for the application to listen on
EXPOSE 8080

# Run the application when the container starts
CMD ["./web"]
