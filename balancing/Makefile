.PHONY: init migrate down

init:
	docker-compose up -d --build jaeger auth-storage-master ;\
	docker image prune -f ;\

migrate:
	docker-compose -f docker-compose.yml up --build auth-storage-master ;\
	docker rm -f auth-storage-master ;\
	docker image prune -f ;\

down:
	docker-compose down ;\
