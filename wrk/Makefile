.PHONEY: launch_1 launch_10 launch_100 launch_1000

launch_1:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1Mzk1MGEwNS0xOWRiLTExZWItYjg5Ni0wMjQyYWMxMjAwMDIiLCJleHAiOjE2MDM5ODIwNzEsImlhdCI6MTYwMzk3ODQ3MSwiaXNzIjoiYXV0aF9zZXJ2aWNlIiwibmJmIjoxNjAzOTc4NDcxfQ.rVbiB7joH3w-6nN_vNtA4mOVkn2kOk-U6PsSW72JhNQ" \
		-t 1 \
		-c 1 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_10:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1Mzk1MGEwNS0xOWRiLTExZWItYjg5Ni0wMjQyYWMxMjAwMDIiLCJleHAiOjE2MDM5ODIwNzEsImlhdCI6MTYwMzk3ODQ3MSwiaXNzIjoiYXV0aF9zZXJ2aWNlIiwibmJmIjoxNjAzOTc4NDcxfQ.rVbiB7joH3w-6nN_vNtA4mOVkn2kOk-U6PsSW72JhNQ" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_100:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1Mzk1MGEwNS0xOWRiLTExZWItYjg5Ni0wMjQyYWMxMjAwMDIiLCJleHAiOjE2MDM5ODIwNzEsImlhdCI6MTYwMzk3ODQ3MSwiaXNzIjoiYXV0aF9zZXJ2aWNlIiwibmJmIjoxNjAzOTc4NDcxfQ.rVbiB7joH3w-6nN_vNtA4mOVkn2kOk-U6PsSW72JhNQ" \
		-t 10 \
		-c 100 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_1000:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1Mzk1MGEwNS0xOWRiLTExZWItYjg5Ni0wMjQyYWMxMjAwMDIiLCJleHAiOjE2MDM5ODIwNzEsImlhdCI6MTYwMzk3ODQ3MSwiaXNzIjoiYXV0aF9zZXJ2aWNlIiwibmJmIjoxNjAzOTc4NDcxfQ.rVbiB7joH3w-6nN_vNtA4mOVkn2kOk-U6PsSW72JhNQ" \
		-t 1000 \
		-c 1000 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\
