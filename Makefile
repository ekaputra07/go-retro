npm-deps:
	cd web/assets && npm install

npm-build:
	cd web/assets && npm run build

go-deps:
	go mod tidy

setup: go-deps npm-deps

build:
	./scripts/build

test:
	./scripts/test

compose: npm-build
	docker compose up