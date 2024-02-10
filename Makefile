all: d-build-api-gateway d-compose

d-build-api-gateway: .
	docker build -t dc-api-gateway:local ./api-gateway

d-compose:
	docker compose up