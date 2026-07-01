# TaskManagementSystemAPI

## 📖 Project overview
A task management application consisting of:
- Go REST API (Gin)
- CodeIgniter 4 Web Client
- MariaDB
- JWT Authentication with Refresh Tokens

## ✨ Features
- User Authentication
- JWT Access Tokens
- Refresh Tokens
- Task CRUD
- File Attachments
- User Profile
- Dashboard
- Search
- Pagination
- Role-Based Access Control

## 🛠️ Technology stack
Backend API
- Go
- Gin
- MariaDB

Frontend
- CodeIgniter 4
- Bootstrap 5

Authentication
- JWT
- Refresh Tokens

## 📁 Project architecture (CI4 → Go API → MariaDB)
+-----------------------+
|   CodeIgniter 4 Web   |
| (Presentation Layer)  |
+-----------+-----------+
            |
            | HTTP Requests
            | JWT Access Token
            v
+-----------------------+
|     Go REST API       |
|   (Business Layer)    |
+-----------+-----------+
            |
            | SQL Queries
            v
+-----------------------+
|        MariaDB        |
|      (Data Layer)     |
+-----------------------+

**CodeIgniter 4 Web Application**
The CodeIgniter 4 application serves as the user interface of the system. It is responsible for:
- User interface and page rendering
- Form validation
- Session management
- Storing Access Tokens and Refresh Tokens
- Communicating with the Go REST API
- Automatically refreshing expired JWT Access Tokens
- Displaying API responses to users

The CI4 application does not access the database directly. All business operations are performed through the REST API.

**Go REST API**
The Go API is the core backend of the system and contains all business logic. Responsibilities include:
- User authentication
- JWT Access Token generation
- Refresh Token management
- Authorization
- Task management
- User management
- Dashboard statistics
- File upload and attachment management
- Profile picture management
- Database operations

Every request from the web application is validated and processed by the API before interacting with the database.

**MariaDB Database**
MariaDB stores all application data, including:
- Users
- Roles
- Tasks
- Statuses
- Tags
- File attachment metadata
- Refresh Tokens

The database is accessed exclusively through the Go API.

## ⚙️ Installation guide
1. Clone the repository
2. Install dependencies
3. Create the database