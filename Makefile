ui:
	./scripts/npm-dev

setup:
	./scripts/go-install
	./scripts/npm-install

test:
	./scripts/test

compose:
	./scripts/npm-build
	docker compose up