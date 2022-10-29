run:
	go run cmd/posterr/main.go

env-up:
	docker-compose up --build -d db

env-down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

cover:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html
	open coverage.html

generate-query:
	sqlc generate

run-docker:
	docker-compose up --build -d app

seed:
	go run cmd/seed/main.go

fmt:
	gofmt -s -w .