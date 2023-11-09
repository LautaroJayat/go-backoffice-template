FROM golang:alpine as builder
ARG DOCKER_USER=app_user
RUN addgroup -S $DOCKER_USER && adduser -S $DOCKER_USER -G $DOCKER_USER
USER $DOCKER_USER
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir out
RUN go build -o=out/main cmd/server/main.go 

FROM alpine:latest
# remember to pass a user in --build-arg
ARG DOCKER_USER=app_user
RUN addgroup -S $DOCKER_USER && adduser -S $DOCKER_USER -G $DOCKER_USER
USER $DOCKER_USER
WORKDIR /app
COPY --from=builder /app/out/main /app/main
COPY --from=builder /app/config/default.yaml /app/config/default.yaml
CMD ["/app/main"]