
.PHONY: run
run:
	go run cmd/http/main.go

.PHONY: build
build:
	docker-compose up -d --build