FROM golang:1.22-alpine AS build

WORKDIR /go/src/distributed-calculator/page-parser

COPY app ./app/

COPY go.mod ./

RUN go mod download

RUN go build -o ../../../bin/app ./app/main.go

FROM alpine
WORKDIR /go

COPY web/app/build ./web/app/build/

COPY --from=build /go/bin/app /bin/app

CMD ["app"]