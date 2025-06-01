# Pin Go version for reproducible builds
FROM golang:1.24.1-alpine AS build

WORKDIR /app

# download dependencies - low frequency of changes
COPY go.mod go.sum ./
RUN go mod download

# copy sources - higher frequency of changes
COPY internal/ internal/
COPY cmd/ cmd/
COPY api/ api/
COPY scripts/ scripts/

# ensure tests are passing
RUN go test ./...

# create executable - ambulance-webapi-srv
# we want to use scratch image so setting
# the build options in the way that will link all dependencies statically
RUN CGO_ENABLED=0 GOOS=linux \
      go build \
      -ldflags="-w -s" \
      -installsuffix 'static' \
      -o ./ambulance-webapi-srv ./cmd/api

############################################ 

FROM scratch

WORKDIR /app

# see https://github.com/opencontainers/image-spec/blob/main/annotations.md for details
LABEL org.opencontainers.image.authors="Roman Osadsky"
LABEL org.opencontainers.image.title="Ambulance Waiting List WebAPI Service"
LABEL org.opencontainers.image.description="WEBAPI for managing entries in ambulances` waiting list"

# list all variables and their default values for clarity
ENV AMBULANCE_API_ENVIRONMENT=production
ENV AMBULANCE_API_PORT=8080
ENV AMBULANCE_API_MONGODB_HOST=mongo
ENV AMBULANCE_API_MONGODB_PORT=27017
ENV AMBULANCE_API_MONGODB_DATABASE=hospital-spaces
ENV AMBULANCE_API_MONGODB_USERNAME=root
ENV AMBULANCE_API_MONGODB_TIMEOUT_SECONDS=5

COPY --from=build /app/ambulance-webapi-srv ./

# Actual port may be changed during runtime
# Default using for the simple case scenario
EXPOSE 8080
ENTRYPOINT ["./ambulance-webapi-srv"]