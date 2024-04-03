all: d-build-auth d-build-api-gateway d-build-daemon d-build-page-parser d-compose

d-build-api-gateway: .
	docker build -t dc-api-gateway:local ./api-gateway

d-build-auth: .
	docker build -t dc-auth:local ./auth

d-build-daemon: .
	docker build -t dc-daemon:local ./daemon

d-build-page-parser: .
	docker build -t dc-page-parser:local ./page-parser

d-compose:
	docker compose up