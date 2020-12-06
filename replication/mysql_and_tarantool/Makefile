.PHONY: init migrate replicator backend down

init:
	docker-compose up -d storage_master storage_tarantool ;\

migrate:
	docker-compose -f docker-compose.yml up --build migrator ;\
	docker image prune -f ;\
	docker rm -f mysql-migrator-tarantool ;\

replicator:
	docker-compose -f docker-compose.yml up --build -d replicator ;\
	docker image prune -f ;\

backend:
	docker-compose -f docker-compose.yml up --build -d backend ;\
	docker image prune -f ;\

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
