# Use Go's official image to build the app
FROM golang:1.24-alpine AS builder

# Set the current working directory in the container
WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Generate Swagger docs (this will create the docs folder)
RUN swag init --generalInfo cmd/server/main.go --output docs

# Build the Go app
RUN go build -o main cmd/server/main.go

# Use a lightweight image to run the app
FROM alpine:3.15

# Install necessary libraries (e.g., CA certificates for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary and the docs folder into the new image
COPY --from=builder /app/main .
COPY --from=builder /app/docs /root/docs

# Expose the port your app will run on
EXPOSE 8080

# Run the application
CMD ["./main"]
