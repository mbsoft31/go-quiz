# Stage 1: Build the Go app
FROM golang:1.20 as builder

WORKDIR /app

# Install Templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the Go modules manifest and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /quiz-app

# Stage 2: Run the app in a lightweight container
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built app binary from the builder stage
COPY --from=builder /quiz-app .

# Copy static files
COPY --from=builder /app/public /public

# Expose the application port
EXPOSE 4000

# Run the app
CMD ["./quiz-app"]
