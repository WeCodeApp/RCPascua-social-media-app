# Task Management API

A RESTful API for a simple task management application with Microsoft/Outlook authentication.

## Features

- User authentication with Microsoft/Outlook accounts
- JWT-based authentication for API endpoints
- Task management (create, read, update, delete)
- MySQL database integration with GORM
- Database migrations and seeding with faker data
- Proper logging
- MVC architecture

## Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Microsoft Azure AD application (for authentication)

## Configuration

Create a `.env` file in the root directory with the following variables:

```
PORT=8080
JWT_SECRET=your-secret-key
MICROSOFT_CLIENT_ID=your-microsoft-client-id
MICROSOFT_CLIENT_SECRET=your-microsoft-client-secret
MICROSOFT_REDIRECT_URI=http://localhost:8080/auth/microsoft/callback
MICROSOFT_TENANT_ID=common

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-db-password
DB_NAME=go_azure

# Application Environment (development, production)
APP_ENV=development
```

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Create a MySQL database:
   ```sql
   CREATE DATABASE go_azure CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
4. Configure your `.env` file with the database credentials
5. Goto the cmd/migration folder and Run the application: `go run main.go`
   - The application will automatically run migrations and seed the database in development mode
   - To disable seeding, set `APP_ENV=production` in your `.env` file

## API Endpoints

### Authentication

- `GET /auth/microsoft`: Initiates Microsoft OAuth login
- `GET /auth/microsoft/callback`: Handles the callback from Microsoft OAuth
- `POST /auth/signout`: Signs out the user

### Tasks

All task endpoints require authentication with a JWT token in the Authorization header.

- `GET /tasks`: Get all tasks for the authenticated user
- `GET /tasks/:id`: Get a specific task by ID
- `POST /tasks`: Create a new task
- `PUT /tasks/:id`: Update an existing task
- `DELETE /tasks/:id`: Delete a task

## Authentication Flow

1. The client redirects the user to `/auth/microsoft`
2. The user logs in with their Microsoft/Outlook account
3. Microsoft redirects back to `/auth/microsoft/callback` with an authorization code
4. The server exchanges the code for a token and returns a JWT token to the client
5. The client includes the JWT token in the Authorization header for subsequent requests

## Task Model

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "completed": false,
  "user_id": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Example Requests

### Create a Task

```
POST /tasks
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Complete project",
  "description": "Finish the task management API project",
  "completed": false
}
```

### Get All Tasks

```
GET /tasks
Authorization: Bearer <jwt_token>
```

### Update a Task

```
PUT /tasks/:id
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Complete project",
  "description": "Finish the task management API project",
  "completed": true
}
```

### Delete a Task

```
DELETE /tasks/:id
Authorization: Bearer <jwt_token>
```

## License

MIT
