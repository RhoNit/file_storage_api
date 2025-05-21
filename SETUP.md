### Used Modules / Libraries

- [`uber-go/zap`](https://github.com/uber-go/zap) – Used as a logger
- [`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt/v5) – Used for JWT utilities
- [`joho/godotenv`](https://github.com/joho/godotenv) – Used as an environment variable loader
- [`labstack/echo/v4`](https://github.com/labstack/echo/v4) – Used as the Echo web server framework


### Install required tools

#### Install Swagger
```bash
# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest
```

#### Install jq (optional, for better JSON formatting)
```bash
# For macOS
brew install jq

# For Ubuntu/Debian
sudo apt-get install jq
```

## Project Setup

### 1. Clone the repository
```bash
git clone https://github.com/RhoNit/file_storage_api.git
cd file_storage_api
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Configure environment
```bash
# Create .env file from example
# although in hurriedness, i didn't put things on .env file.
# But vars like `JWT_SECRET`, `MAX_ALLOWED_STORAGE_PER_USER`, `DEFAULT_PAGE_SIZE` etc could be kept in the `.env` file 
```

### 4. Generate Swagger documentation
```bash
swag init .
```

### 5. Build and Run

```bash
# Build the application
go build -o file-storage-api main.go

# Run the application
./file-storage-api

# Or directly run the application
go run main.go
```
