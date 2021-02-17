.PHONY: init migrate app down

init:
	docker-compose -f docker-compose.yml up -d --build ch-cluster ch-shard-0 ch-shard-1 auth-storage \
    		social-storage cache nats-streaming jaeger consul-server ;\
	docker image prune -f ;\

migrate:
	docker-compose -f docker-compose.yml up --build auth-migrator social-migrator ch-cluster-migrator \
		ch-shard-migrator-0 ch-shard-migrator-1 ;\
	docker rm -f auth-migrator social-migrator ch-cluster-migrator ch-migrator-0 ch-migrator-1 ;\
	docker image prune -f ;\

app:
	docker-compose -f docker-compose.yml up -d --build auth social messenger gateway ;\
	docker image prune -f ;\

down:
	docker-compose down ;\
