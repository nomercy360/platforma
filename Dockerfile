# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /app

COPY . .

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /api main.go

FROM alpine:3.19 AS build-release-stage

RUN apk --no-cache add ca-certificates bash curl

WORKDIR /app

COPY --from=build-stage /api /app/api

EXPOSE 8080

CMD ["/api"]
