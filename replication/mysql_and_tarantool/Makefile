.PHONY: init migrate replicator backend down

init:
	docker-compose up -d storage_master storage_tarantool ;\

migrate:
	docker-compose -f docker-compose.yml up --build migrator ;\
	docker rm -f mysql_migrator_replica ;\

replicator:
	docker-compose -f docker-compose.yml up --build -d replicator ;\

backend:
	docker-compose -f docker-compose.yml up --build -d backend ;\

wrk:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=ma ;\

down:
	docker-compose down ;\
