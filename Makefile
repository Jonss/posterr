run-local:
	go run cmd/posterr/main.go

env-up:
	docker-compose up --build -d db

run-docker:
	docker-compose up --build app

run: env-up run-docker

env-down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

cover:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html
	open coverage.html

gen-query:
	sqlc generate

seed:
	go run cmd/seed/main.go

fmt:
	gofmt -s -w .

gen-mock:
	mockgen -destination db/mock/db.go github.com/Jonss/posterr/db AppQuerier
	mockgen -destination pkg/post/mock/service.go github.com/Jonss/posterr/pkg/post Service
