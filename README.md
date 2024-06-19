### Library Management System Backend with GO lang

Welcome to the Library Management System Backend! This RESTful API is designed from the perspective of a librarian, providing an efficient way to manage users, books, and issued books. Built while learning Go, this API showcases practical usage of the Gin framework and MySQL for database operations.
Before going futher this is the ERD of the system

![Untitled](https://github.com/gandharvtalikoti/lib-management-system/assets/79464855/c787b4a4-b1c3-44f2-a47d-a6c97a7a3a92)


## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Create User](#create-user)
  - [Get Users](#get-users)
  - [Get Specific User](#get-specific-user)
  - [Issue Book](#issue-book)
  - [Get Books Issued to a User](#get-books-issued-to-a-user)
  - [Get Overdue Books for a User](#get-overdue-books-for-a-user)
  - [Get All Overdue Books](#get-all-overdue-books)
  - [Search for a Book](#search-for-a-book)
  - [Get User Details](#get-user-details)
  - [Get Book Log](#get-book-log)
- [API Endpoints](#api-endpoints)

## Features

- Create, read, update, and delete users
- Issue books to users
- Track books issued to users
- Monitor overdue books
- Search for books by name
- View detailed user and book logs

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/library-management-system.git
   cd library-management-system
   ```

2. **Install dependencies**:
   Make sure you have Go installed. Then, run:
   ```bash
   go mod tidy
   ```

3. **Set up MySQL database**:
   Ensure you have a MySQL database running. Create a database named `lib-sys` and import the necessary schema.

4. **Configure database connection**:
   Update the `database.go` file with your MySQL connection string:
   ```go
   const dbURL = "mysql://root:1910@localhost:3306/lib-sys"
   ```

5. **Run the server**:
   ```bash
   go run main.go
   ```

## Usage

### Create User

- **Endpoint**: `POST /users`
- **Description**: Creates a new user.
- **Example**:
  ```bash
  curl -X POST http://localhost:8080/users -d '{"name":"John Doe", "email":"john@example.com"}'
  ```

### Get Users

- **Endpoint**: `GET /users`
- **Description**: Retrieves a list of all users.
- **Example**:
  ```bash
  curl http://localhost:8080/users
  ```

### Get Specific User

- **Endpoint**: `GET /users/:id`
- **Description**: Retrieves details of a specific user.
- **Example**:
  ```bash
  curl http://localhost:8080/users/1
  ```

### Issue Book

- **Endpoint**: `POST /issued`
- **Description**: Issues a book to a user.
- **Request Body**:
  ```json
  {
      "user_id": 3,
      "book_id": 4
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:8080/issued -d '{"user_id":3, "book_id":4}'
  ```

### Get Books Issued to a User

- **Endpoint**: `GET /issued/:id`
- **Description**: Retrieves a list of books issued to a user.
- **Example**:
  ```bash
  curl http://localhost:8080/issued/3
  ```

### Get Overdue Books for a User

- **Endpoint**: `GET /overdue/:id`
- **Description**: Retrieves a list of overdue books for a user.
- **Example**:
  ```bash
  curl http://localhost:8080/overdue/2
  ```

### Get All Overdue Books

- **Endpoint**: `GET /issued/overdue`
- **Description**: Retrieves a list of all overdue books.
- **Example**:
  ```bash
  curl http://localhost:8080/issued/overdue
  ```

### Search for a Book

- **Endpoint**: `GET /search`
- **Description**: Searches for books by name.
- **Query Parameter**: `name`
- **Example**:
  ```bash
  curl http://localhost:8080/search?name=ikigai
  ```

### Get User Details

- **Endpoint**: `GET /users/details/:user_id`
- **Description**: Retrieves user details along with a list of books they have issued and their overdue statuses.
- **Example**:
  ```bash
  curl http://localhost:8080/users/details/3
  ```

### Get Book Log

- **Endpoint**: `GET /books/book/log/:book_id`
- **Description**: Provides information about a specific book and a list of users who have issued it.
- **Example**:
  ```bash
  curl http://localhost:8080/books/book/log/4
  ```

## API Endpoints

| Method | Endpoint                           | Description                                                       |
|--------|------------------------------------|-------------------------------------------------------------------|
| POST   | /users                             | Creates a new user                                                |
| GET    | /users                             | Retrieves a list of all users                                     |
| GET    | /users/:id                         | Retrieves details of a specific user                              |
| POST   | /issued                            | Issues a book to a user                                           |
| GET    | /issued/:id                        | Retrieves a list of books issued to a user                        |
| GET    | /overdue/:id                       | Retrieves a list of overdue books for a user                      |
| GET    | /issued/overdue                    | Retrieves a list of all overdue books                             |
| GET    | /search                            | Searches for books by name                                        |
| GET    | /users/details/:user_id            | Retrieves user details and a list of books they have issued       |
| GET    | /books/book/log/:book_id           | Provides information about a specific book and its issue history  |

## Conclusion

This Library Management System Backend provides a comprehensive set of features to manage library operations effectively. Feel free to explore and contribute to the project!


### Contact

For any questions or suggestions, please contact (mailto:gandharvwork@example.com).

