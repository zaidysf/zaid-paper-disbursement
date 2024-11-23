FROM golang:1.21-alpine

# Install required system dependencies
RUN apk add --no-cache make gcc musl-dev git

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files if they exist
COPY go.* ./

# Initialize go module
RUN go mod init zaid-paper-disbursement && \
    go mod tidy

# Copy everything from the current directory to the container
COPY . .

# Download dependencies again after copying all source code
RUN go mod tidy

# Install specific version of migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run migrations, seed, and start the app
CMD sh -c "sleep 5 && \
    migrate -path migrations -database postgres://postgres:postgres@db:5432/disbursement_db?sslmode=disable up && \
    go run seeds/seed.go && \
    ./main"