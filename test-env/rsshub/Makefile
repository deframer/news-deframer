.PHONY: start stop down logs clean

start:
	docker compose up -d --build --remove-orphans

stop:
	docker compose stop

down:
	docker compose down

logs:
	docker compose logs -f

clean:
	docker compose down -v