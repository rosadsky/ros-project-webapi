param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:AMBULANCE_API_PORT="8080"
$env:AMBULANCE_API_MONGODB_USERNAME="root"
$env:AMBULANCE_API_MONGODB_PASSWORD="neUhaDnes"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "openapi" {
        docker run --rm -ti  -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd/api
        } finally {
            mongo down
        }
    }
    "mongo" {
        mongo up
    }
    "test" {
        go test -v ./...
    }
    "docker" {
        docker build -t xosadsky/ros-project-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
    }
    default {
        throw "Unknown command: $command"
    }
}