.PHONY: infrastructure migrate service reload_auth reload_social reload_messenger reload_gateway reload_frontend down

infrastructure:
	docker-compose -f deployment/docker-compose.yml up -d --build ch-cluster ch-shard-0 ch-shard-1 auth-storage \
		social-storage cache nats-streaming jaeger ;\

migrate:
	docker-compose -f deployment/docker-compose.yml up --build auth-migrator social-migrator ch-cluster-migrator \
		ch-shard-migrator-0 ch-shard-migrator-1 ;\
	docker rm -f auth-migrator social-migrator ch-cluster-migrator ch-migrator-0 ch-migrator-1 ;\
	docker image prune -f ;\

service:
	docker-compose -f deployment/docker-compose.yml up -d --build auth social messenger gateway ;\
	docker image prune -f ;\

reload_auth:
	docker rm -f auth-otus ;\
	docker-compose -f deployment/docker-compose.yml up -d --build auth ;\
    docker image prune -f ;\

reload_social:
	docker rm -f social-otus ;\
	docker-compose -f deployment/docker-compose.yml up -d --build social ;\
    docker image prune -f ;\

reload_messenger:
	docker rm -f messenger-otus ;\
	docker-compose -f deployment/docker-compose.yml up -d --build messenger ;\
    docker image prune -f ;\

reload_gateway:
	docker rm -f gateway-otus ;\
	docker-compose -f deployment/docker-compose.yml up -d --build gateway ;\
    docker image prune -f ;\

reload_frontend:
	docker rm -f frontend ;\
	docker-compose -f deployment/docker-compose.yml up -d --build frontend ;\
    docker image prune -f ;\

down:
	docker-compose -f deployment/docker-compose.yml down ;\
