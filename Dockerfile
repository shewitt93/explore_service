FROM golang:1.24.0-alpine AS builder

WORKDIR /app


# install dependencies
COPY go.mod go.sum ./
RUN go mod download


# copy source code
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

COPY main.go .
COPY .env .env


RUN go build -o /app/main

