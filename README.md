# Paper.id Technical Assessment

This repository contains my solution for the Paper.id technical assessment.

## Assessment Task
Create a disbursement API with the following requirements:
- User has a balance in the application wallet
- Balance can be disbursed through an API endpoint
- Implemented in Golang
- Single endpoint for disbursement only
- User data and balances can be stored as hard coded or database (I chose to use PostgreSQL for better data persistence and transaction handling)

## Solution Overview

I've implemented a RESTful API service that:
- Manages user wallet balances
- Provides a disbursement endpoint
- Uses PostgreSQL for data storage
- Implements proper transaction handling
- Includes comprehensive tests
- Uses Docker for easy setup and deployment

## Running the Application

Prerequisites:
- Docker
- Docker Compose

Setup and Run:
```bash
# Clone the repository
git clone https://github.com/zaidysf/zaid-paper-disbursement.git
cd zaid-paper-disbursement

# Copy environment file
cp .env.example .env

# Make the database script executable
chmod +x scripts/create-multiple-postgresql-databases.sh

# Start the application
docker-compose up --build
```

## Testing the API

Once the application is running, you can test the disbursement endpoint:

```bash
# Test disbursement endpoint (success case)
curl -X POST http://localhost:8080/api/v1/disbursement \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 100.00
  }'

# Test disbursement endpoint (insufficient balance)
curl -X POST http://localhost:8080/api/v1/disbursement \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 2000.00
  }'
```

<img width="671" alt="image" src="https://github.com/user-attachments/assets/3aa5d1a1-230f-4aeb-8637-a694b25a4d60">


## Running Tests

```bash
# Run all tests
docker-compose exec app go test ./tests/... -v
```

<img width="538" alt="image" src="https://github.com/user-attachments/assets/ad97b9fb-3cbf-42e7-a044-a9331df3ab41">


## Project Structure
```
zaid-paper-disbursement/
├── api/
│   ├── handlers/      # HTTP handlers
│   ├── middlewares/   # HTTP middlewares
│   └── routes/        # Route definitions
├── config/           # Database configuration
├── internal/
│   ├── models/       # Data models
│   └── services/     # Business logic
├── migrations/       # Database migrations
├── scripts/         # Utility scripts
├── seeds/           # Database seeders
├── tests/           # Integration tests
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Author
Zaid Yasyaf (zaid.ug@gmail.com)

## Notes
- The implementation uses PostgreSQL instead of hardcoded data for better reliability and transaction support
- Test coverage includes both success and failure scenarios
- Docker setup ensures consistent environment across different machines
