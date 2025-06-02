FROM golang:1.19 AS build
WORKDIR /go/src
COPY internal/ros-project ./internal/ros-project
COPY main.go .
COPY go.sum .
COPY go.mod .

ENV CGO_ENABLED=0

RUN go build -o ros-project .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/ros-project ./
EXPOSE 8080/tcp
ENTRYPOINT ["./ros-project"]
