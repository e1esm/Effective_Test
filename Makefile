dockerize:
	docker-compose up --build
test:
	go test ./...