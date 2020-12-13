.PHONY: init migrate down

init:
	docker-compose up -d --build ch-cluster ch-shard-0 ch-shard-1 ch-shard-2 ch-shard-3 ch-shard-4 mysql-storage redis \
		backend ;\
	docker image prune -f ;\

migrate:
	docker-compose -f docker-compose.yml up --build mysql-migrator ch-cluster-migrator ch-shard-migrator-0 \
		ch-shard-migrator-1 ch-shard-migrator-2 ch-shard-migrator-3 ch-shard-migrator-4 ;\
    docker rm -f mysql-migrator ch-cluster-migrator ch-migrator-0 ch-migrator-1 ch-migrator-2 ch-migrator-3 \
    	ch-shard-migrator-4 ;\
	docker image prune -f ;\

down:
	docker-compose down ;\
