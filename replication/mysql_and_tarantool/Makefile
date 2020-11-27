.PHONY: init migrate replicator down

init:
	docker-compose up -d storage_mysql storage_tarantool ;\

migrate:
	docker-compose -f docker-compose.yml up --build migrator  ;\
	docker rm -f mysql_migrator_replica ;\

replicator:
	docker-compose -f docker-compose.yml up --build replicator  ;\

down:
	docker-compose down ;\
