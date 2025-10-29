npm-deps:
	cd web/assets && npm install

npm-build:
	cd web/assets && npm run build

go-deps:
	go mod tidy

setup: go-deps npm-deps

dev: npm-build
	GORETRO_HOST=localhost \
	GORETRO_SESSION_SECRET=Bve8zfg8RvNJHh8jxxEAVj8oe00bE2QY \
	GORETRO_SESSION_SECURE=false \
	go run .

build:
	go build -v ./...

test:
	go test -v ./...