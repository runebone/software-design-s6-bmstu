SCRIPT := scripts/generate_env.py

all: up

build:
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml build
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.mongo.yml build

up: postgres

postgres:
	python ${SCRIPT}
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml up
	# cat .env | grep -v '^$$' | sed 's/^\(\w\)/export \1/g' | sed 's/$$/;/g' | xargs && docker compose

mongo:
	python ${SCRIPT}
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.mongo.yml up

userpg:
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml up user-postgres

rmv:
	for volume in $$(docker volume ls -q); do docker volume rm $$volume; done

down:
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.postgres.yml down
	docker compose --env-file .env -f docker-compose.yml -f docker-compose.user.mongo.yml down

boards:
	docker exec -it todo-postgres psql -U postgres -d todo_db -c "select * from boards;"

columns:
	docker exec -it todo-postgres psql -U postgres -d todo_db -c "select * from columns;"

cards:
	docker exec -it todo-postgres psql -U postgres -d todo_db -c "select * from cards;"

users:
	docker exec -it user-postgres psql -U postgres -d user_db -c "select * from users;"

tokens:
	docker exec -it auth-postgres psql -U postgres -d auth_db -c "select * from tokens;"
