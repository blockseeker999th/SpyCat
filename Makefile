.SILENT:

build:
	docker compose up --build

restart:
	docker compose down
	docker compose up

rebuild:
	docker compose down
	docker compose up --build