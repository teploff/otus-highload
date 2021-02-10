.PHONY: init migrate app down

init:
	docker-compose up -d --build jaeger auth-storage-master ;\
	docker image prune -f ;\

migrate:
	docker-compose -f docker-compose.yml up --build auth-migrator ;\
	docker rm -f auth-migrator ;\
	docker image prune -f ;\

app:
	docker-compose -f docker-compose.yml up -d --build nginx ;\
	docker image prune -f ;\

wrk:
	echo ${ACCESS_TOKEN} ;\
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost/auth/user/get-by-email?email=teploff@email.com ;\

down:
	docker-compose down ;\
