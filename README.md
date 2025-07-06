# Tickitz Backend API

A full-featured cinema ticket booking RESTful API built with Go (Gin), PostgreSQL, and JWT authentication. This project supports user registration, movie management, ticket purchasing, profile editing, and role-based access control for admins. It uses a modular architecture, Redis caching, and Swagger documentation.

## Features
- User registration, login, profile edit, and password reset
- Role-based access: admin & regular user
- Admin movie management (create, update, delete, assign genres/directors/casts)
- View all movies, upcoming, and now showing (with search + Redis cache)
- Payment method creation (admin)
- Transaction flow: book tickets with movie, time, seat, and payment method
- JWT-based authentication & authorization
- Swagger documentation ready


## Tech Stack
- Go (Gin) for backend framework
- PostgreSQL with pgx driver
- Redis for caching GET requests
- JWT for secure authentication
- Swagger (via swaggo/swag) for API docs
- Docker-ready (optional)

## How to Clone and Use

Make sure you have Golang installed on your device.

#### 1. Clone the repository
```
git clone https://github.com/ranandasatria/fgo24-be-tickitz.git
```

#### 2. Navigate into the project directory
```
cd fgo24-be-tickitz
```

#### 3. Install the dependencies
```
go mod tidy
```

#### 4. Setup .env 
Create a .env file in the root folder with the following variables:
```
APP_SECRET=your_jwt_secret
DATABASE_URL=postgres://username:password@localhost:5433/tontrix-db
RDADDRESS=localhost:6379
RDPASSWORD=
RDDB=0
```

#### 5. Run the program
```
go run main.go
```

## Authentication
Most endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your_token_here>
```

## API Endpoints

| Method | Endpoint             | Description                        | Auth Required |
|--------|----------------------|------------------------------------|---------------|
AUTH
| POST | /register | Register a new user | ❌ |
| POST | /login | Login and receive JWT token | ❌ |
| POST | /forgot-password | Send password reset token via email | ❌ |
| POST | /reset-password | Reset password using token | ❌ |
USERS
| GET | /users | Get all users | ✅ admin |
| DELETE | /admin/users/{id}| Delete a user  | ✅ admin |
PROFILE
| GET | /profile | Get logged-in user profile | ✅ |
| PATCH | /profile | Edit profile and optionally password| ✅ |
MOVIES
| GET | /movies | List all movies (with search & pagination) | ❌ |
| GET | /movies/{id} | Get movie details by ID | ❌ |
| GET | /movies/now-showing | Get now showing movies | ❌ |
| GET | /movies/upcoming | Get upcoming movies | ❌ |
| POST | /admin/movies | Create new movie | ✅ admin |
| PATCH | /admin/movies/{id} | Update movie details | ✅ admin |
| DELETE | /admin/movies/{id} | Delete a movie | ✅ admin |
Genres, Directors, Actors
| GET | /genres | Get all genres | ❌ |
| POST | /admin/genres | Create a genre | ✅ admin |
| DELETE | /admin/genres/{id} | Delete a genre | ✅ admin |
| GET | /directors | Get all directors | ❌ |
| POST | /admin/directors | Create a director | ✅ admin |
| DELETE | /admin/directors/{id}| Delete a director | ✅ admin |
| GET | /actors | Get all actors | ❌ |
| POST | /admin/actors | Create an actor | ✅ admin |
| DELETE | /admin/actors/{id} | Delete an actor | ✅ admin |
Payment Methods
| GET | /admin/payment-method | View all payment methods | ✅ admin |
| POST | /admin/payment-method | Add a new payment method | ✅ admin |
| DELETE | /admin/payment-method/{id} | Delete a payment method | ✅ admin |
 Transactions
| GET | /transactions | Get logged-in user's transactions | ✅ |
| POST | /transactions | Create a new transaction | ✅ |
| GET | /admin/transactions | View all transactions  | ✅ admin |


# ENTITY-RELATIONSHIP DIAGRAM 

```mermaid
erDiagram
direction LR

users ||--o{ movies : adds
users ||--o{ transactions : books
movies ||--o{ transactions : has
genres ||--o{ movie_genres : categorize
movie_genres }o--|| movies : has
directors ||--o{ movie_directors : directs
movie_directors }o--|| movies : has
actors ||--o{ movie_casts : plays
movie_casts }o--|| movies : has
transactions ||--o{ transaction_details : has
transactions }o--|| payment_method : used

users {
  int id PK
  varchar email
  varchar password
  varchar full_name
  varchar phone_number
  varchar profile_picture
  varchar role
  timestamp created_at
  timestamp updated_at
}

movies {
  int id PK
  varchar title
  text description
  date release_date
  int duration_minutes
  varchar image
  varchar horizontal_image
  timestamp created_at
  timestamp updated_at
}

genres {
  int id PK
  varchar genre_name
  timestamp created_at
  timestamp updated_at
}

movie_genres {
  int id PK
  int id_movie FK
  int id_genre FK
  timestamp created_at
  timestamp updated_at
}

directors {
  int id PK
  varchar director_name
  timestamp created_at
  timestamp updated_at
}

movie_directors {
  int id PK
  int id_movie FK
  int id_director FK
  timestamp created_at
  timestamp updated_at
}

actors {
  int id PK
  varchar actor_name
  timestamp created_at
  timestamp updated_at
}

movie_casts {
  int id PK
  int id_movie FK
  int id_actor FK
  timestamp created_at
  timestamp updated_at
}

transactions {
  int id PK
  int id_user FK
  int id_movie FK
  date show_date
  time show_time
  varchar cinema
  varchar location
  int total_price
  int payment_method FK
  timestamp created_at
  timestamp updated_at
}

transaction_details {
  int id PK
  int transaction_id FK
  varchar seat
  int price
  timestamp created_at
  timestamp updated_at
}

payment_method {
  int id PK
  varchar payment_name
  timestamp created_at
  timestamp updated_at
}

```

## 📄 License

This project is licensed under the **MIT License**.  

## ©️ Copyright

&copy; 2025 Kodacademy

