run:
	go run cmd/posterr/main.go

env-up:
	docker-compose up --build -d db

env-down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

cover:
	go tool cover -html=coverage.out

generate-query:
	sqlc generate

run-docker:
	docker-compose up --build -d app

seed:
	go run cmd/seed/main.go

