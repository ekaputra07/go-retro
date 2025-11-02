npm-deps:
	cd ui/assets && npm install

npm-build:
	cd ui/assets && npm run build

go-deps:
	go mod tidy

setup: go-deps npm-deps

dev: npm-build
	go run ./cmd/web -secret=Bve8zfg8RvNJHh8jxxEAVj8oe00bE2QY
build:
	go build -v ./cmd/web -o go-retro

test:
	go test -v ./...