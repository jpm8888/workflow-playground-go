# Workflow Project

This project implements a workflow engine that manages various workflows, including authentication and order processing. It is built using Go and utilizes GORM for database interactions.

## Project Structure

- **cmd/api**: Contains the main application entry point.
- **internal/api**: Contains HTTP handlers for authentication and order management.
- **internal/auth**: Contains models, services, and workflows related to user authentication.
- **internal/order**: Contains models, services, and workflows related to order processing.
- **pkg/workflow**: Contains the core workflow engine, including types, errors, and engine logic.
- **go.mod**: Go module file for dependency management.

## Features

- User authentication workflow with OTP verification and PIN setup.
- Order processing workflow with order creation, payment confirmation, and order fulfillment.
- GORM integration for database operations.

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   ```

2. Navigate to the project directory:
   ```
   cd app/myproj
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Run the application:
   ```
   go run cmd/api/main.go
   ```

5. Access the API at `http://localhost:8080`.

## API Endpoints

### Authentication

- `POST /auth/start`: Start the authentication process.
- `POST /auth/verify-otp`: Verify the OTP sent to the user's phone.
- `POST /auth/set-pin`: Set the user's PIN.
- `POST /auth/verify-pin`: Verify the user's PIN.

### Orders

- `POST /orders`: Create a new order.
- `POST /orders/:id/payment`: Confirm payment for an order.
- `POST /orders/:id/fulfill`: Fulfill an order.

## License

This project is licensed under the MIT License. See the LICENSE file for details.