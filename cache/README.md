# Отчет о домашнем задании №5. Кеширование

## Содержание
1. [ Задание ](#task)
    - [ Цель ](#task-goal)
    - [ Приобретенные навыки ](#task-skills)
    - [ Постановка задачи ](#task-statement)
2. [ Сведения ](#information)
    - [ Используемые инструменты ](#information-tools)
    - [ Характеристики железа ](#information-computer)
3. [ Ход работы ](#work)
4. [ Итоги ](#results)

<img align="right" width="480" src="static/title-page.png">

<a name="task"></a>
## Задание
Кеширование ленты новостей социальной сети.

<a name="task-goal"></a>
### Цель
Реализовать механизм кеширования ленты новостей социальной сети.

<a name="task-skills"></a>
### Приобретенные навыки
В результате выполненного домашнего задания необходимо приобрести следующие навыки:
- работа с кешами;
- работа с очередями;
- проектирование масштабируемых архитектур.

<a name="task-statement"></a>
### Постановка задачи
Реализовать функционал просмотра новостей (своих и друзей) в отдельной web-странице.
Для этого необходимо:
- реализовать утверждение/отклонение заявок в друзья пользователями социальной сети;
- удаление из списка друзей;
- добавление и просмотр новостей.

Держать последние новости необходимо в горячем кеше и обновлять его через очередь, чтоб сгладить нагрузку при  пиковых 
запросах. То есть обновление идет асинхронное, но мы всегда читаем ленту из кеша. Предусмотреть механизм перестройки
кешей: инвалидация и прогрев.

<a name="information"></a>
## Сведения
<a name="information-tools"></a>
### Используемые инструменты
Для выполнения дз понадобятся следующие инструменты:
- [golang](https://golang.org/dl/) (version 1.14 or 1.15)
- [docker](https://docs.docker.com/get-docker/) (>= version 19.03.8) & [docker compose](https://docs.docker.com/compose/install/) (>= version 1.25.5);
- [jq](https://stedolan.github.io/jq/download/) (>= version 1.5)

<a name="information-computer"></a>
### Характеристики железа
Домашнее задание выполнялось на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC

<a name="work"></a>
## Ход работы
В ходе выполнения домашнего задания в качестве кеша выступает [Redis](https://redis.io/), а в качетве очереди 
[Nats-Streaming](https://docs.nats.io/nats-streaming-concepts/intro).

Склонируем наш проект:
```shell
git clone https://github.com/teploff/otus-highload.git && cd otus-highload
```

Поднимем всю инфраструктуру и применим миграции:
```shell
make up && make migrate
```

Далее необходимо зарегистрировать пользоватейлей. Это можно сделать с помощью WebUI, расположенного по адресу: 
http://localhost:8081/sign-up или же curl-ом. В данном случае для эффективности воспольюсь именно curl-ом.

Зарегистрируем трех пользователей в системе:
- Шерлока Холмса,
- Доктора Ватсона,
- Джеймса Мориарти.

```shell script
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "holmes@gmail.com", "password": "12345678", "name": "Sherlock", "surname": "Holmes", "birthday": "1854-01-09T20:21:25+00:00", "sex": "male", "city": "London", "interests": "the violin, smoking, investigation"}' \
    http://localhost:9999/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "watson@gmail.com", "password": "12345678", "name": "John", "surname": "Watson", "birthday": "1852-08-7T20:21:25+00:00", "sex": "male", "city": "London", "interests": "to heal people, investigation"}' \
    http://localhost:9999/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "moriarty@gmail.com", "password": "12345678", "name": "James", "surname": "Moriarty", "birthday": "1835-05-01T20:21:25+00:00", "sex": "male", "city": "London", "interests": "the crime"}' \
    http://localhost:9999/auth/sign-up
```

Получим access token'ы от системы для наших пользователей и запишем их в переменные окружения:
```shell script
export HOLMES_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
-d '{"email": "holmes@gmail.com", "password": "1234567890"}' \
http://localhost:9999/auth/sign-in | jq -r '.access_token')
export WATSON_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
-d '{"email": "watson@gmail.com", "password": "1234567890"}' \
http://localhost:9999/auth/sign-in | jq -r '.access_token')
export MORIARTY_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
-d '{"email": "moriarty@gmail.com", "password": "1234567890"}' \
http://localhost:9999/auth/sign-in | jq -r '.access_token')
```

Проверим наличие access token-ов:
```shell script
echo $HOLMES_ACCESS_TOKEN
echo $WATSON_ACCESS_TOKEN
echo $MORIARTY_ACCESS_TOKEN
```

Получим ID-шники созданных пользователей, они нам пригодятся для отправки заявок в друзья и их подтверждении.
```shell script
export HOLMES_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${WATSON_ACCESS_TOKEN}" \
    http://localhost:9999/auth/user?email=holmes@gmail.com | jq -r '.user_id')
export WATSON_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
  http://localhost:9999/auth/user?email=holmes@gmail.com | jq -r '.user_id')
export MORIARTY_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
  http://localhost:9999/auth/user?email=holmes@gmail.com | jq -r '.user_id') 
```

Проверим, что запрос успешно выполнился, применив команду:
```shell script
echo $HOLMES_ID
echo $WATSON_ID
echo $MORIARTY_ID
```

Теперь давайте со стороны Доктора Ватсона и Джеймса Мориарти отправим заявки в друзья Шерлоку Холмсу:
```shell script
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${WATSON_ACCESS_TOKEN}" \
    -d '{
         "friends_id": ["'"$HOLMES_ID"'"]
        }' \
    http://localhost:9999/social/create-friendship
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${MORIARTY_ACCESS_TOKEN}" \
    -d '{
         "friends_id": ["'"$HOLMES_ID"'"]
        }' \
    http://localhost:9999/social/create-friendship
```

Проверим, что у Шерлока Холмса появилось 2 follower-а:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
    http://localhost:9999/social/followers
```

Подтвердим заявки в друзья, сделав follower-ов своими друзьями:
```shell script
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
    -d '{
         "friends_id": ["'"$WATSON_ID"'", "'"$MORIARTY_ID"'"]
        }' \
    http://localhost:9999/social/confirm-friendship
```

Проверим то, что у Шерлока Холмса появилось два новых друга:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
    http://localhost:9999/social/friends
```

Теперь необходимо проверить функционал новостей.
Под новостью понимаем маленькое сообщение или твит о том, что нового у пользователя:)
Создадим пару новостей от Доктора Ватсона и пару новостей Джеймса Мориарти, также одну новость от Шерлока Холмса.
```shell script
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${WATSON_ACCESS_TOKEN}" \
    -d '{
         "news": ["I'm getting married!"]
        }' \
    http://localhost:9999/social/create-news
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${WATSON_ACCESS_TOKEN}" \
    -d '{
         "news": ["Merry is a pretty woman!"]
        }' \
    http://localhost:9999/social/create-news
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${MORIARTY_ACCESS_TOKEN}" \
    -d '{
         "news": ["I'm a nightmare!", "Holmes, you are lost!"]
        }' \
    http://localhost:9999/social/create-news
curl -X POST -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
    -d '{
         "news": ["Greate news!"]
        }' \
    http://localhost:9999/social/create-news
```

Теперь запросим все доступные новости для трех наших персонажей.

Доктором Ватсоном:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${WATSON_ACCESS_TOKEN}" \
    http://localhost:9999/social/news
```

Должны получить нечто следующее:</br>
<p align="center">
    <img src="static/algorithm_chose_shard.png">
</p>

Джеймсом Мориарти:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${MORIARTY_ACCESS_TOKEN}" \
    http://localhost:9999/social/news
```

Должны получить нечто следующее:</br>
<p align="center">
    <img src="static/algorithm_chose_shard.png">
</p>

И Шерлоком Холмсом:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${HOLMES_ACCESS_TOKEN}" \
    http://localhost:9999/social/news
```

Должны получить нечто следующее:</br>
<p align="center">
    <img src="static/algorithm_chose_shard.png">
</p>

Технические моменты
Инвалидация кеша осуществляется по событию. Какое-либо изменение в базе данных влечет за собой валидацию, а именно:
- добавление, удалиление друзей;
- добавление новостей.

Прогрев кеша осуществляется в двух случаях:
- при старте системы (запускается джоба на прогрев кеша и только после успеха система запускается);
- при получение события из очереди (если по каким-то причинам, мы понимаем, что кеша невалидный запускаем его прогрев).

Давайте проверим это.

Первый способ

Зайдем в Redis и полностью сброим кеш:
```shell script
docker exec -it redis-cache redis-cli
AUTH password
FLUSHALL
exit
```

Теперь пересоберем backend:
```shell script
 make reload_backend
```

После того, как backend пересобрался, перейдем в Redis и проверим BD 1 и 2 (1 - Друзья, 2 - Новости).
```shell script
docker exec -it redis-cache redis-cli
AUTH password
SELECT 1
KEYS *
SELECT 2
KEYS *
```

Должны получить нечто следующее:</br>
<p align="center">
    <img src="static/algorithm_chose_shard.png">
</p>
Видим, что кеш при старте системы прогревается.

Второй способ

Зайдем в Redis и полностью сброим кеш:
```shell script
docker exec -it redis-cache redis-cli
AUTH password
FLUSHALL
exit
```

Теперь не трогая backend, запустим простенькую tool'зу, которая в очередь шлет уведомление, на основе которого
система понимает, что необходимо заново прогресть кеш:
```shell script
cd backend/tools/cache-heater-enabler
go run main.go --addr="localhost:4222" --cluster_id="stan-cluster" --subject="cache-reload"
```

Перейдем в Redis и проверим наличие кеша:
```shell script
docker exec -it redis-cache redis-cli
AUTH password
SELECT 1
KEYS *
SELECT 2
KEYS *
```


docker exec -it redis-cache redis-cli
AUTH password
FLUSHALL

```shell
cd backend/tools/cache-heater-enabler
go run main.go --addr="localhost:4222" --cluster_id="stan-cluster" --subject="cache-reload"
```

Должны получить нечто следующее:</br>
<p align="center">
    <img src="static/algorithm_chose_shard.png">
</p>
Видим, что кеш при отправке сообщения в очередь прогревается.

<a name="results"></a>
## Итоги
В ходе выполнения домашнего задания: