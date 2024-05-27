gen:
	go run github.com/99designs/gqlgen generate

docker-build:
	docker build -t ozon_journal .

res:
	docker compose down
	docker compose up --build