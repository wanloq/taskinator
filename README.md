#  Taskinator

Taskinator is a task management system built using the **Fiber framework**, **PostgreSQL**, and **Golang best practices**. This project serves as an exploration of backend development, authentication, database management, and API design in Golang.

## 📌 Features

- ✅ User Authentication (JWT-based)
- ✅ Role-based Access Control (Admin/User)
- ✅ Task Creation, Updating, and Deletion
- ✅ PostgreSQL Database Integration (via GORM)
- ✅ API Documentation with Swagger
- ✅ Secure Password Hashing (bcrypt)
- ✅ Middleware for Authentication & Authorization
- ✅ Structured and Scalable Codebase
- ✅ Docker Support (for containerized deployment)

## 📂 Project Structure

```

📂 taskinator/
│── 📂 .github/workflows/               # GitHub Actions CI/CD workflow
│   └── go.yml                           # Go CI pipeline
│
│── 📂 db/migrations/                     # Database migrations (SQL files)
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_add_updated_at_to_users.up.sql
│   ├── 000002_add_updated_at_to_users.down.sql
│   ├── 000003_add_deleted_at_to_users.up.sql
│   ├── 000003_add_deleted_at_to_users.down.sql
│   ├── 000004_add_password_at_to_users.up.sql
│   ├── 000004_add_password_at_to_users.down.sql
│
│── 📂 docs/                              # API Documentation (Swagger, Postman, etc.)
│
│── 📂 internal/                           # Core application code
│   │── 📂 config/                        # Configurations and database setup
│   │   ├── .air.toml                     # Live reload config for Air
│   │   ├── config.go                     # App configuration handling
│   │   ├── db.go                         # Database connection logic
│   │   ├── migrations.go                 # Migration handling
│   │
│   │── 📂 controllers/                   # Route handlers (business logic)
│   │   ├── user_controller.go            # User-related logic
│   │
│   │── 📂 dto/                           # Data Transfer Objects (DTOs)
│   │   ├── auth_dto.go                   # DTOs for authentication
│   │
│   │── 📂 middleware/                    # Middleware for authentication, logging, etc.
│   │   ├── auth_middleware.go            # Authentication middleware
│   │
│   │── 📂 models/                        # Database models
│   │   ├── user.go                       # User model definition
│   │
│   │── 📂 repositories/                  # Database query logic
│   │   ├── user_repository.go            # User data access logic
│   │
│   │── 📂 routes/                        # API route definitions
│   │   ├── routes.go                     # Main route registry
│   │   ├── user_routes.go                # User-specific routes
│   │
│   │── 📂 utils/                         # Utility functions
│   │   ├── jwt.go                        # JWT token handling
│   │   ├── password.go                   # Password hashing and validation
│
│── 📂 task-manager-frontend/              # Frontend (if applicable)
│── 📂 tmp/                                # Temporary files
│── .env                                  # Environment variables (DO NOT COMMIT!)
│── .env.example                          # Example environment file
│── .gitignore                            # Git ignore rules
│── go.mod                                # Go module dependencies
│── go.sum                                # Go module checksums
│── main.go                               # Main application entry point
│── README.md                             # Project documentation

```

## 🛠️ Setup & Installation

### 1️⃣ Clone the repository
```sh
git clone https://github.com/yourusername/taskinator.git
cd taskinator
```

### 2️⃣ Configure Environment Variables
Create a `.env` file and specify your database credentials:
See the  `.env.example` file for complete content.

### 3️⃣ Install Dependencies
```sh
go mod tidy
```

### 4️⃣ Run Database Migrations
```sh
make migrate-up   # Apply database migrations
```

### 5️⃣ Start the Application
```sh
go run main.go
```

### 6️⃣ API Documentation (Swagger)
Once the server is running, open your browser and visit:
```
http://localhost:3000/swagger/index.html
```

## 🚀 API Endpoints

| Method  | Endpoint       | Description                   | Auth Required |
|---------|---------------|-------------------------------|--------------|
| `POST`  | `/register`   | Register a new user          | ❌ No |
| `POST`  | `/login`      | Authenticate user & get JWT  | ❌ No |
| `GET`   | `/tasks`      | Get all tasks (user-specific) | ✅ Yes |
| `POST`  | `/tasks`      | Create a new task            | ✅ Yes |
| `PUT`   | `/tasks/:id`  | Update a task                | ✅ Yes |
| `DELETE`| `/tasks/:id`  | Delete a task                | ✅ Yes |

## 🐳 Docker (Optional)
To run Taskinator in a Docker container, use:
```sh
docker-compose up --build
```

## 🤝 Contributing
Contributions are welcome! Feel free to fork the repository, submit issues, or create pull requests.

---

### 📜 License
This project is licensed under the **MIT License**.

---

## 📌 Notes
- This project follows **Golang best practices** (layered architecture).
- Uses **golang-migrate** instead of `AutoMigrate` for production readiness.
- Designed to be scalable and easily extendable.

---