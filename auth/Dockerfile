FROM golang:1.22-alpine AS build

WORKDIR /go/src/distributed-calculator/auth

COPY app ./app/

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o ../../../bin/app ./app/cmd/app/main.go

FROM alpine
WORKDIR /go

COPY ./configs ./configs
COPY --from=build /go/bin/app ./app