🔐 CloudCart — Auth Service

1. Purpose

The Auth Service is responsible for:

- User registration
- Password hashing
- User login
- JWT access-token generation
- JWT validation
- Authentication middleware
- Refresh-token generation and storage
- Generating new access tokens using refresh tokens

2. Technology Used

```
Go
│
├── Gin              → HTTP framework
├── GORM             → ORM
├── PostgreSQL       → Database
├── bcrypt           → Password hashing
└── golang-jwt/jwt   → JWT authentication
```

3. Service Structure
```
services/
└── auth-service/
    │
    ├── cmd/
    │   └── server/
    │       └── main.go
    │
    ├── internal/
    │   │
    │   ├── app/
    │   │   └── app.go
    │   │
    │   ├── config/
    │   │   └── config.go
    │   │
    │   ├── database/
    │   │   ├── database.go
    │   │   ├── migrate.go
    │   │   └── seed.go
    │   │
    │   ├── logger/
    │   │   └── logger.go
    │   │
    │   ├── model/
    │   │   ├── user.go
    │   │   ├── role.go
    │   │   └── refresh_token.go
    │   │
    │   ├── repository/
    │   │   ├── user_repository.go
    │   │   └── user_repository_impl.go
    │   │
    │   ├── service/
    │   │   ├── auth_service.go
    │   │   └── dto.go
    │   │
    │   ├── handler/
    │   │   └── auth_handler.go
    │   │
    │   ├── router/
    │   │   └── router.go
    │   │
    │   ├── auth/
    │   │   ├── jwt.go
    │   │   └── refresh_token.go
    │   │
    │   └── middleware/
    │       └── auth_middleware.go
    │
    ├── .env
    ├── .gitignore
    └── go.mod
```

4. Application Startup Flow
We deliberately kept main.go simple.
```
main.go
   │
   ▼
app.New()
   │
   ├── Load configuration
   ├── Initialize logger
   ├── Connect PostgreSQL
   ├── Create auth schema
   ├── AutoMigrate tables
   ├── Seed CUSTOMER role
   ├── Create JWT Manager
   ├── Create Repository
   ├── Create Service
   └── Create Handler
   │
   ▼
app.Run()
   │
   ▼
Gin Router
   │
   ▼
HTTP Server
```

5. Configuration
Sensitive/configurable values come from environment variables.
```
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=xxxxx
DB_NAME=cloudcart

JWT_SECRET=cloudcart-development-secret
JWT_EXPIRATION=24h
```

config.go reads these values and creates:
```
type Config struct {
    // DB configuration...

    JWTSecret     string
    JWTExpiration time.Duration
}
```

Important concept
We don't hardcode the JWT secret inside jwt.go.
```
.env
 ↓
config.go
 ↓
JWTManager
 ↓
GenerateToken()
ValidateToken()
```
Later, when deploying to Kubernetes/AWS, this can be replaced with Kubernetes Secrets/AWS Secrets Manager without changing JWT logic.

6. PostgreSQL Design
**One PostgreSQL database + multiple schemas
For Auth:
```
cloudcart
   │
   └── auth
       ├── roles
       ├── users
       └── refresh_tokens
```
This was a deliberate learning/project decision rather than creating separate RDS instances for every microservice.

7. Database Models
auth.roles
```
id
name
created_at
updated_at
```

Currently seeded:
```
CUSTOMER
```

auth.users
```
id
first_name
last_name
email
password_hash
role_id
created_at
updated_at
```
Relationship:
```
users.role_id
      │
      ▼
roles.id
```

auth.refresh_tokens
```
id
user_id
token
expires_at
created_at
updated_at
```
Relationship:
```
refresh_tokens.user_id
          │
          ▼
       users.id
```

8. GORM Schema Handling
Each model has:
```
func (User) TableName() string {
    return "auth.users"
}
```
Similarly:
```
func (Role) TableName() string {
    return "auth.roles"
}
```
and:
```
func (RefreshToken) TableName() string {
    return "auth.refresh_tokens"
}
```
GORM automatically detects TableName().
We also create the schema:
```
CREATE SCHEMA IF NOT EXISTS auth;
```
Then:
```
db.AutoMigrate(...)
```
creates/updates the tables.

9. Repository Layer
The repository abstracts database operations away from business logic.
Interface:
```
type UserRepository interface {
    Create(user *model.User) error
    FindByEmail(email string) (*model.User, error)
    FindRoleByName(name string) (*model.Role, error)

    CreateRefreshToken(token *model.RefreshToken) error
    FindRefreshToken(token string) (*model.RefreshToken, error)
}
```
Important pattern:
```
Handler
   ↓
Service
   ↓
Repository
   ↓
GORM
   ↓
PostgreSQL
```
The service doesn't directly execute SQL.

10. Registration Flow
Endpoint:
```POST /api/v1/auth/register```

Request:
```
{
  "first_name": "Vidhya",
  "last_name": "Shekar",
  "email": "user@example.com",
  "password": "password"
}
```
Flow:
```
Request
  ↓
Gin Handler
  ↓
Auth Service
  ↓
Check existing email
  ↓
Find CUSTOMER role
  ↓
bcrypt password hashing
  ↓
Create User
  ↓
PostgreSQL
  ↓
201 Created
```
Password is never stored as plain text.
Database contains something like:
```$2a$10$................```

11. Login Flow
Endpoint:
```POST /api/v1/auth/login```

Request:
```
{
  "email": "user@example.com",
  "password": "password"
}
```

Flow:
```
Email + Password
       ↓
Find User
       ↓
bcrypt.CompareHashAndPassword()
       ↓
      Valid?
       │
       ▼
Generate Access JWT
       +
Generate Refresh Token
       │
       ▼
Store Refresh Token in DB
       │
       ▼
Return both tokens
```

Response:
```
{
  "access_token": "eyJ...",
  "refresh_token": "....",
  "token_type": "Bearer"
}
```
Wrong credentials return:
```401 Unauthorized```
with:
```
{
  "error": "invalid email or password"
}
```

12. JWT Access Token
JWT contains claims such as:
```
user_id
email
role
exp
iat
```
Conceptually:
```
JWT
 ├── user_id
 ├── email
 ├── role
 ├── issued-at
 └── expiration
 ```
 Signing:
 ```
 HS256
+
JWT_SECRET
```
The JWT manager has two primary responsibilities:
```
GenerateToken()
ValidateToken()
```

13. Authentication Middleware
File:
```internal/middleware/auth_middleware.go```
The middleware:

1. Reads Authorization header
2. Checks Bearer <token>
3. Validates JWT signature
4. Checks expiration
5. Extracts claims
6. Stores them in Gin context
7. Allows request to continue

Request:
```Authorization: Bearer eyJ...```

Flow:
```
Request
   ↓
Auth Middleware
   ↓
Authorization Header
   ↓
Extract JWT
   ↓
Validate JWT
   ↓
Extract Claims
   ↓
Gin Context
   ↓
Protected Handler
```

Claims are stored as:
```
c.Set("user_id", claims.UserID)
c.Set("email", claims.Email)
c.Set("role", claims.Role)
```

14. /me Protected API
Endpoint:
```GET /api/v1/auth/me ```

This endpoint demonstrates that the JWT middleware actually works.
Without token:

```401 Unauthorized```

With:

```Authorization: Bearer <access-token>```

response:
```
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "CUSTOMER"
}
```
The important point is that /me gets this information from the validated JWT, rather than asking the client to send the user information.

15. Refresh Token
Access token and refresh token have different purposes.
```
Access Token
     │
     └── Used for API authorization
         Short-lived


Refresh Token
     │
     └── Used to obtain new access token
         Longer-lived
         Stored in DB
```

Current refresh-token lifetime:
```7 days```
The refresh token is generated using cryptographically secure random bytes.

16. Refresh API
Endpoint:
```POST /api/v1/auth/refresh```

Request:
```
{
  "refresh_token": "..."
}
```

Flow:
```
Refresh Token
      ↓
Find token in DB
      ↓
Check expiration
      ↓
Get associated user
      ↓
Generate new Access JWT
      ↓
Return new Access Token
```
Response:
```
{
  "access_token": "eyJ...",
  "token_type": "Bearer"
}
```
This endpoint is not protected by the access-token middleware because its purpose is to obtain a new access token when the existing one may have expired.

17. Final Auth API List
```
| API                     | Method | Auth   |
| ----------------------- | ------ | ------ |
| `/health`               | GET    | Public |
| `/api/v1/auth/register` | POST   | Public |
| `/api/v1/auth/login`    | POST   | Public |
| `/api/v1/auth/refresh`  | POST   | Public |
| `/api/v1/auth/me`       | GET    | 🔒 JWT |

```

18. Complete Authentication Architecture
                    ┌──────────────┐
                    │    Client    │
                    └──────┬───────┘
                           │
             ┌─────────────┼─────────────┐
             │             │             │
             ▼             ▼             ▼
         Register        Login         /me
             │             │             │
             │             ▼             │
             │       Validate user      │
             │             │             │
             │       ┌─────┴─────┐       │
             │       ▼           ▼       │
             │   Access JWT   Refresh    │
             │                    │       │
             │                    ▼       │
             │             PostgreSQL     │
             │                           │
             │                           ▼
             │                    JWT Middleware
             │                           │
             │                    Validate JWT
             │                           │
             │                           ▼
             │                       /me Handler
             │
             └───────────────┬───────────┘
                             ▼
                         PostgreSQL


19. What You Have Practically Learned
The Auth Service has covered several important Go/backend concepts:
```
Go
├── Package organization
├── Interfaces
├── Dependency injection
├── Structs
├── Error handling
└── time.Duration

Gin
├── Routing
├── Handlers
├── JSON binding
├── Middleware
└── Context

GORM
├── Models
├── Relationships
├── Preload
├── AutoMigrate
└── PostgreSQL schemas

Security
├── bcrypt
├── JWT
├── Access tokens
├── Refresh tokens
├── Authorization header
└── Middleware

Architecture
├── Handler
├── Service
├── Repository
├── Database
└── Configuration
```