npm-install:
	cd web/ui && npm install

npm-dev:
	cd web/ui && npm run dev

go-deps:
	go mod tidy

setup: go-deps npm-install

build:
	./scripts/build

build-ui:
	./scripts/build-ui
	
test:
	./scripts/test

compose: build-ui
	docker compose up