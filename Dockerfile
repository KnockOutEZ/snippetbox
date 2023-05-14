# Start from a Golang v1.16 base image
FROM golang:1.20

# Set the working directory to /app
WORKDIR /app

# Copy the necessary files and directories to the container
COPY . .
# Install the necessary dependencies
# RUN go get github.com/gorilla/mux
# RUN go get github.com/lib/pq

# Build the Go application
RUN go build -o main ./cmd/web

# Expose port 8080
EXPOSE 4000

# Set the environment variables
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=mysecretpassword
ENV DB_NAME=mydb

# Start the Go application
CMD ["./main"]