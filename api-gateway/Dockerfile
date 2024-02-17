FROM golang:1.22-alpine as build

WORKDIR /go/src/distributed-calculator/api-gateway

COPY go.mod go.sum ./

RUN go mod download

COPY app ./app/

RUN go build -o ../../../bin/app ./app/cmd/app/main.go

FROM alpine
WORKDIR /go

COPY configs ./configs/
COPY ./.env .

COPY --from=build /go/bin/app /bin/app

CMD ["app"]