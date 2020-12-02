.PHONY: init migrate down

init:
	docker-compose up -d db-node-1 db-node-2 db-node-3 db-node-4 proxysql redis ;\
	docker-compose up -d --build backend ;\

migrate:
	docker-compose -f docker-compose.yml up --build migrator-1 ;\
    docker rm -f mysql-migrator-sharding-1 ;\
    docker-compose -f docker-compose.yml up --build migrator-2 ;\
	docker rm -f mysql-migrator-sharding-2 ;\
	docker-compose -f docker-compose.yml up --build migrator-3 ;\
	docker rm -f mysql-migrator-sharding-3 ;\
	docker-compose -f docker-compose.yml up --build migrator-4 ;\
	docker rm -f mysql-migrator-sharding-4 ;\

down:
	docker-compose down ;\
