npm-deps:
	cd web/assets && npm install

npm-build:
	cd web/assets && npm run build

go-deps:
	go mod tidy

setup: go-deps npm-deps

build:
	go build -v -o dist/goretro-web ./cmd/web

test:
	go test -v ./...

run:
	go run ./cmd/web -secret dev_Bve8zfg8RvNJHh8jxxEAVj8oe00bE2QY

dev: npm-build run
