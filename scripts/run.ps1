param(
    [Parameter(Mandatory = $true)]
    [string]$command
)

$ProjectRoot = $PSScriptRoot + "/.."

# Set MongoDB environment variables
$env:AMBULANCE_API_PORT = "8080"
$env:AMBULANCE_API_MONGODB_USERNAME = "root"
$env:AMBULANCE_API_MONGODB_PASSWORD = "neUhaDnes"
$env:AMBULANCE_API_MONGODB_DATABASE = "hospital-spaces"
$env:AMBULANCE_API_MONGODB_URI = "mongodb://${env:AMBULANCE_API_MONGODB_USERNAME}:${env:AMBULANCE_API_MONGODB_PASSWORD}@localhost:27017/${env:AMBULANCE_API_MONGODB_DATABASE}?authSource=admin"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "start" {
        Write-Host "Starting Hospital Spaces API with MongoDB..." -ForegroundColor Green
        try {
            Write-Host "Starting MongoDB containers..." -ForegroundColor Yellow
            mongo up --detach
            
            Write-Host "Waiting for MongoDB to be ready..." -ForegroundColor Yellow
            Start-Sleep -Seconds 10
            
            Write-Host "Environment variables set:" -ForegroundColor Yellow
            Write-Host "  AMBULANCE_API_MONGODB_URI: $env:AMBULANCE_API_MONGODB_URI" -ForegroundColor Gray
            Write-Host "  AMBULANCE_API_MONGODB_DATABASE: $env:AMBULANCE_API_MONGODB_DATABASE" -ForegroundColor Gray
            Write-Host "  AMBULANCE_API_PORT: $env:AMBULANCE_API_PORT" -ForegroundColor Gray
            Write-Host "  MongoDB Web UI: http://localhost:8081 (mexpress/mexpress)" -ForegroundColor Gray
            
            Write-Host "Starting Go application..." -ForegroundColor Green
            go run ${ProjectRoot}/cmd/api
        } finally {
            Write-Host "Stopping MongoDB containers..." -ForegroundColor Yellow
            mongo down
        }
    }
    "build" {
        Write-Host "Building Hospital Spaces API..." -ForegroundColor Green
        go build -o ${ProjectRoot}/bin/server ${ProjectRoot}/cmd/api
        Write-Host "Build completed: bin/server" -ForegroundColor Green
    }
    "dev" {
        Write-Host "Starting development server with MongoDB..." -ForegroundColor Green
        try {
            Write-Host "Starting MongoDB containers..." -ForegroundColor Yellow
            mongo up --detach
            
            Write-Host "Waiting for MongoDB to be ready..." -ForegroundColor Yellow
            Start-Sleep -Seconds 10
            
            Write-Host "Environment variables set:" -ForegroundColor Yellow
            Write-Host "  AMBULANCE_API_MONGODB_URI: $env:AMBULANCE_API_MONGODB_URI" -ForegroundColor Gray
            Write-Host "  AMBULANCE_API_MONGODB_DATABASE: $env:AMBULANCE_API_MONGODB_DATABASE" -ForegroundColor Gray
            Write-Host "  AMBULANCE_API_PORT: $env:AMBULANCE_API_PORT" -ForegroundColor Gray
            Write-Host "  MongoDB Web UI: http://localhost:8081 (mexpress/mexpress)" -ForegroundColor Gray
            
            Write-Host "Starting development server with hot reload..." -ForegroundColor Green
            air
        } finally {
            Write-Host "Stopping MongoDB containers..." -ForegroundColor Yellow
            mongo down
        }
    }
    "mongo" {
        Write-Host "Starting MongoDB containers..." -ForegroundColor Green
        Write-Host "MongoDB will be available at: mongodb://localhost:27017" -ForegroundColor Yellow
        Write-Host "MongoDB Web UI will be available at: http://localhost:8081" -ForegroundColor Yellow
        Write-Host "Web UI credentials: mexpress/mexpress" -ForegroundColor Gray
        Write-Host "Database credentials: ${env:AMBULANCE_API_MONGODB_USERNAME}/${env:AMBULANCE_API_MONGODB_PASSWORD}" -ForegroundColor Gray
        mongo up
    }
    "mongo-down" {
        Write-Host "Stopping MongoDB containers..." -ForegroundColor Yellow
        mongo down
    }
    "test" {
        Write-Host "Running tests..." -ForegroundColor Green
        go test ./...
    }
    "docker" {
        docker run -e AMBULANCE_API_MONGODB_PASSWORD=neUhaDnes hospital-spaces-api:local-build
    }
    "clean" {
        Write-Host "Cleaning up..." -ForegroundColor Yellow
        if (Test-Path "${ProjectRoot}/bin") {
            Remove-Item -Recurse -Force "${ProjectRoot}/bin"
            Write-Host "Removed bin directory" -ForegroundColor Yellow
        }
        if (Test-Path "${ProjectRoot}/tmp") {
            Remove-Item -Recurse -Force "${ProjectRoot}/tmp"
            Write-Host "Removed tmp directory" -ForegroundColor Yellow
        }
    }
    default {
        Write-Host "Hospital Spaces Management Script" -ForegroundColor Green
        Write-Host ""
        Write-Host "Prerequisites:" -ForegroundColor Yellow
        Write-Host "  - Go 1.21 or later" -ForegroundColor White
        Write-Host "  - Docker and Docker Compose" -ForegroundColor White
        Write-Host "  - Air for development (install with: go install github.com/cosmtrek/air@latest)" -ForegroundColor White
        Write-Host ""
        Write-Host "Available commands:" -ForegroundColor Yellow
        Write-Host "  start       - Run the application with MongoDB" -ForegroundColor White
        Write-Host "  dev         - Run in development mode with hot reload and MongoDB" -ForegroundColor White
        Write-Host "  build       - Build the application" -ForegroundColor White
        Write-Host "  test        - Run tests" -ForegroundColor White
        Write-Host "  mongo       - Start MongoDB containers only" -ForegroundColor White
        Write-Host "  mongo-down  - Stop MongoDB containers" -ForegroundColor White
        Write-Host "  clean       - Clean up build artifacts" -ForegroundColor White
        Write-Host ""
        Write-Host "MongoDB Configuration:" -ForegroundColor Yellow
        Write-Host "  Database: ${env:AMBULANCE_API_MONGODB_DATABASE}" -ForegroundColor Gray
        Write-Host "  Username: ${env:AMBULANCE_API_MONGODB_USERNAME}" -ForegroundColor Gray
        Write-Host "  Password: ${env:AMBULANCE_API_MONGODB_PASSWORD}" -ForegroundColor Gray
        Write-Host "  URI: ${env:AMBULANCE_API_MONGODB_URI}" -ForegroundColor Gray
        Write-Host "  Web UI: http://localhost:8081 (mexpress/mexpress)" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Examples:" -ForegroundColor Yellow
        Write-Host "  ./scripts/run.ps1 dev" -ForegroundColor Gray
        Write-Host "  ./scripts/run.ps1 mongo" -ForegroundColor Gray
        Write-Host "  ./scripts/run.ps1 start" -ForegroundColor Gray
    }
} 