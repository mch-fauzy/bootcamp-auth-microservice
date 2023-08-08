# Bootcamp Auth Microservice

The **Bootcamp Auth Microservice** provides user authentication and authorization functionalities for the bootcamp ecosystem.

## Features

- Student registration
- Login
- Update name by user
- Read user
- JWT-based authentication
- Role-based authorization

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mucha-fauzy/bootcamp-auth-microservice.git
   ```

2. Navigate to the project directory:

   ```bash
   cd bootcamp-auth-microservice
   ```

3. Install the dependencies:

   ```bash
   go mod download
   ```
4. Create Required table (./migrations) and seed(./seeders) the table if necessary

5. Build the project:

   ```bash
   go run main.go
   ```


By default, the microservice will listen on port 8080.

## API Endpoints

- **POST /v1/register**: Register a new user for Student. Requires `username`, `password`, and optional `name` in the request body.

- **POST /v1/login**: Authenticate and generate a JWT token. Requires `username` and `password` in the request body.

- **GET /v1/validate-auth**: Validate a JWT token and retrieve user information.

- **PUT /v1//users/{id}**: Update name of existing users by ID. Requires `name` in the request body and `teacher` role.

- **GET /v1//users**: Read all users, with filter by name and pagination


## Authentication

The microservice uses JWT (JSON Web Tokens) for authentication. Upon successful login, a JWT token is generated and should be included in the `Authorization` header as `Bearer <token>` for protected routes.

## Future Updates
- **PUT /v1//users/{id}**: Changes to only update current login user
- **GET /v1//users**:  Change to Read current login user profile
- Pisahkan app logic dengan db logic (service: app, repo: db)
- Set secret key untuk JWT di .env dan buat configs nya (users-management-crud-api) -> jadi nanti tinggal dipanggil tidak perlu hard code 

---
Feel free to modify and expand this README.md to include any additional details or specific instructions for running and using your bootcamp-auth microservice.