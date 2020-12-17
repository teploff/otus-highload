.PHONY: up down reload_backend reload_frontend

up:
	docker-compose -f docker-compose.yml up --build -d storage adminer backend frontend ;\
	docker image prune -f ;\

migrate:
	docker-compose up --build migrator ;\
	docker rm -f mysql-migrator ;\
	docker image prune -f ;\

reload_backend:
	docker rm -f social_network_backend ;\
	docker-compose -f docker-compose.yml up -d --build social_network_backend ;\
    docker image prune -f ;\

reload_frontend:
	docker rm -f social_network_frontend ;\
	docker-compose -f docker-compose.yml up -d --build social_network_frontend ;\
    docker image prune -f ;\

down:
	docker-compose -f docker-compose.yml down ;\
