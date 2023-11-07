run: build
	./build/taskmaster

build:
	go build -o build/taskmaster ./cmd/taskmaster

migrate: build-migrate
	./build/migrate

build-migrate:
	go build -o build/migrate ./cmd/migrate
