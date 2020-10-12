.PHONY: up down reload_backend reload_frontend

up:
	docker-compose -f docker-compose.yml up -d storage adminer ;\
	docker-compose -f docker-compose.yml up --build migrator ;\
	docker-compose rm -f migrator ;\
	docker-compose -f docker-compose.yml up -d --build social_network_backend social_network_frontend ;\
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
