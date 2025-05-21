# API Documentation (for API Testing)

## Base URL
N.B. The Application's been deployed in ec2 machine. And the public ip is `43.204.211.156`
```
http://43.204.211.156:8085/api
```

## Authentication
All authenticated endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### User Management

#### Register User
- **Endpoint:** `POST /register`
- **Description:** Registers a new user
- **Authentication:** Not required
- **Request Body:**
  ```json
  {
    "username": "ranit",
    "password": "qwerty1234"
  }
  ```
- **cURL:**
  ```bash
  curl -X POST http://43.204.211.156:8085/api/register \
    -H "Content-Type: application/json" \
    -d '{"username": "ranit", "password": "qwerty1234"}' | jq .
  ```
- **Success Response:**
  ```json
  {
    "message": "User registered successfully"
  }
  ```
- **Error Response (Username exists):**
  ```json
  {
    "error": "Username already exists"
  }
  ```

#### Login
- **Endpoint:** `POST /login`
- **Description:** Authenticates user and returns JWT token
- **Authentication:** Not required
- **Request Body:**
  ```json
  {
    "username": "ranit",
    "password": "qwerty1234"
  }
  ```
- **cURL:**
  ```bash
  curl -X POST http://43.204.211.156:8085/api/login \
    -H "Content-Type: application/json" \
    -d '{"username": "ranit", "password": "qwerty1234"}' | jq .
  ```
- **Response:**
  ```json
  {
    "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc3NjMwMjAsInVzZXJuYW1lIjoicmFuaXRfYmlzd2FzIn0.-CzgcRbOGFM9QbvEzCNTp4h98doYgbRLkOpORuzOJDw"
  }
  ```

### File Management

#### Upload File
- **Endpoint:** `POST /upload`
- **Description:** Uploads a file
- **Authentication:** Required
- **Request:**
  - Content-Type: multipart/form-data
  - Body: file
- **cURL:**
  ```bash
  curl -X POST http://43.204.211.156:8085/api/upload \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc3NjMwMjAsInVzZXJuYW1lIjoicmFuaXRfYmlzd2FzIn0.-CzgcRbOGFM9QbvEzCNTp4h98doYgbRLkOpORuzOJDw" \
    -F "file=@/mnt/c/Users/Ranit/Downloads/dbms.pdf"
  ```
- **Response:**
  ```json
  {
    "message": "File uploaded successfully"
  }
  ```

- **cURL:**
  ```bash
  curl -X POST http://43.204.211.156:8085/api/upload \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc3NjMwMjAsInVzZXJuYW1lIjoicmFuaXRfYmlzd2FzIn0.-CzgcRbOGFM9QbvEzCNTp4h98doYgbRLkOpORuzOJDw" \
    -F "file=@/mnt/c/Users/Ranit/Downloads/ZoomInstallerFull.exe"
  ```
- **Response:**
  ```json
  {
    "error":"Storage quota exceeded"
  }
  ```

#### Get User Files
- **Endpoint:** `GET /files`
- **Description:** Fetches list of uploaded files with pagination
- **Authentication:** Required
- **Query Parameters:**
  - `page` (integer): Page number (default: 1)
  - `pageSize` (integer): Items per page (default: 10)
- **cURL:**
  ```bash
  curl -X GET "http://43.204.211.156:8085/api/files?page=2&pageSize=3" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc3NjMwMjAsInVzZXJuYW1lIjoicmFuaXRfYmlzd2FzIn0.-CzgcRbOGFM9QbvEzCNTp4h98doYgbRLkOpORuzOJDw" | jq .
  ```
- **Response:**
  ```json
  {
    "paginated_response": {
      "data": [
        {
          "filename": "main.go",
          "originalName": "main.go",
          "size": 1856,
          "uploadTime": "2025-05-20T22:31:55.556488+05:30",
          "username": "ranit"
        },
        {
          "filename": "docs.go",
          "originalName": "docs.go",
          "size": 10621,
          "uploadTime": "2025-05-20T22:32:06.565904+05:30",
          "username": "ranit"
        },
        {
          "filename": "swagger.json",
          "originalName": "swagger.json",
          "size": 9954,
          "uploadTime": "2025-05-20T22:32:14.704002+05:30",
          "username": "ranit"
        }
      ],
      "page": 2,
      "pageSize": 3,
      "totalItems": 7,
      "totalPages": 3
    }
  }
  ```

### Storage Management

#### Get Remaining Storage
- **Endpoint:** `GET /storage/remaining`
- **Description:** Returns the user's current storage usage and remaining space
- **Authentication:** Required
- **cURL:**
  ```bash
  curl -X GET "http://43.204.211.156:8085/api/storage/remaining" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDc3NjMwMjAsInVzZXJuYW1lIjoicmFuaXRfYmlzd2FzIn0.-CzgcRbOGFM9QbvEzCNTp4h98doYgbRLkOpORuzOJDw" | jq .
  ```
- **Response:**
  ```json
  {
    "storage info": {
      "totalStorage": 52428800,
      "usedStorage": 28885,
      "remainingStorage": 52399915
    },
    "username": "ranit"
  }
  ```
