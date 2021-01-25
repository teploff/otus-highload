.PHONY: db up migrate reload_backend reload_frontend down

db:
	docker-compose -f deployment/docker-compose.yml up -d storage

up:
	docker-compose -f deployment/docker-compose.yml up --build -d adminer gateway frontend ;\
	docker image prune -f ;\

migrate:
	docker-compose -f deployment/docker-compose.yml up --build migrator ;\
	docker rm -f mysql-migrator ;\
	docker image prune -f ;\

reload_backend:
	docker rm -f gateway-otus ;\
	docker-compose -f deployment/docker-compose.yml up -d --build gateway ;\
    docker image prune -f ;\

reload_frontend:
	docker rm -f frontend ;\
	docker-compose -f deployment/docker-compose.yml up -d --build frontend ;\
    docker image prune -f ;\

down:
	docker-compose -f deployment/docker-compose.yml down ;\
