.PHONY: init migrate down

init:
	docker-compose up -d storage_mysql storage_tarantool replicator ;\

migrate:
	docker-compose -f docker-compose.yml up --build migrator  ;\
	docker rm -f mysql_migrator_replica ;\

down:
	docker-compose down ;\
