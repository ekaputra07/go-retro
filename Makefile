ui-install:
	cd web/ui && npm install

ui-build:
	cd web/ui && npm run build

ui-dev:
	cd web/ui && npm run dev

go-deps:
	go mod tidy

setup: go-deps ui-install

build:
	./scripts/build

test:
	./scripts/test

compose: ui-build
	docker compose up