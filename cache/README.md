# Otus Highload. Social network

Инвалидация по событию
Удаляем данные из кеша при изменении в базе 

docker exec -it redis-cache redis-cli
AUTH password
FLUSHALL

```shell
cd backend/tools/cache-heater-enabler
go run main.go --addr="localhost:4222" --cluster_id="stan-cluster" --subject="cache-reload"
```

<img align="right" width="320" src="static/title-page.png">

It's Otus's homework for **Highload Architect** course, which is a basic social network.
The project consists of two parts: *Backend* and *Frontend*.
- Backend part was written with **Golang** using 1.15 version. 
- Frontend part was written with **VueJS** using 2.6.12 version.

The social network allows:
- SignUp/SignIn/Authorize users in the system;
- Get other user's questioners;

Authorization result represents of two JWT's: access & refresh.

## Launch project 

### Requirements
- Docker & docker-compose

### Up infrastructure
```shell script
make up
```

### Shutting down infrastructure
```shell script
make down
```
