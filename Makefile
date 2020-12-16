
.PHONY: run
run:
	go run cmd/http/main.go -K 01234567890123456789012345678901

.PHONY: build
build:
	docker-compose up -d --build