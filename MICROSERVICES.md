# Microservices

## Содержание
1. [ Задание ](#task)
    - [ Цель ](#task-goal)
    - [ Приобретенные навыки ](#task-skills)
    - [ Постановка задачи ](#task-statement)
2. [ Сведения ](#information)
    - [ Используемые инструменты ](#information-tools)
    - [ Характеристики железа ](#information-computer)
3. [ Ход работы ](#work)
    - [ Предложенный вариант решения ](#work-solution)
    - [ Выполнение ](#work-execute)
        - [ Подготовка ](#work-execute-preparation)
    - [ Технические моменты ](#work-technical-moments)
      - [ Серверная часть ](#work-technical-moments-server)
      - [ Клиентская часть ](#work-technical-moments-client)
4. [ Итоги ](#results)

<img align="right" width="480" src="static/microservices/preview.png">

<a name="task"></a>
## Задание
Декомпозиция монолитной инфраструктуры на микросервисы.

<a name="task-goal"></a>
### Цель
Декомпозировать бизнес-домен монолитной инфраструктуры на отдельные микросервисы.

<a name="task-skills"></a>
### Приобретенные навыки
В результате выполненного задания необходимо приобрести следующие навыки:
- декомпозиции предметной области;
- разделение монолитного приложения;
- работа с HTTP;
- работа с REST API и gRPC;
- проектирование масштабируемых архитектур.

<a name="task-statement"></a>
### Постановка задачи
Декомпозировать имеющую монолитную инфраструктуру на микросервисы, каждый из которых выполняет строго определенную 
бизнес-доменную задачу. Для этого необходимо:
- вынести систему диалогов в отдельный микросервис;
- вынести систему регистрации/авторизации/аутентификации в отдельный сервис;
- осуществить взаимодействие между микросервисами по REST API и gRPC;
- организовать сквозное логирование;
- предусмотреть то, что не все клиенты обновляют приложение быстро и кто-то может ходить через старое API.

<a name="information"></a>
## Сведения
<a name="information-tools"></a>
### Используемые инструменты
Для выполнения задания понадобятся следующие инструменты:
- [docker](https://docs.docker.com/get-docker/) (>= version 19.03.8) & [docker compose](https://docs.docker.com/compose/install/) (>= version 1.25.5);
- [curl](https://curl.haxx.se/download.html) (>= version 7.68.0);
- [websocat](https://github.com/vi/websocat/releases) (>= version 1.6.0);
- [jq](https://stedolan.github.io/jq/download/) (>= version 1.5).

<a name="information-computer"></a>
### Характеристики железа
Задание выполнялось на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC


<a name="work"></a>
## Ход работы

<a name="work-solution"></a>
### Предложенный вариант решения

Разработанная микросервисная инфраструктура имеет следующий вид: </br>
<p align="center">
    <img src="static/microservices/infrastructure-schema.png">
</p>
    

<a name="work-execute"></a>
### Выполнение
Для того, чтобы осуществить вышеописанную задумку нам понадобится следующее:
  - собрать cluster Clickhouse, который будет состоять из одной cluster node'ы и пяти node для шардирования;
  - экземпляр кеша(**Redis**-а) для того, чтобы персистентно отслеживать эффект *Lady Gaga'и* у пользователей;


Склонируем наш проект:
```shell
git clone https://github.com/teploff/otus-highload.git && cd otus-highload
```

Поднимаем инфраструктуру, состоящую из:
- одного экземпляра MySQL;
- шести экземпляров ClickHouse (один - cluster, остальные - shard'ы);
- одного экземпляра Redis;
- одного экземпляра backend'а

и применяем миграции:
```shell
make infrastructure && make migrate && make service
```

<a name="work-execute-preparation"></a>
#### Подготовка
Для демонстрации работы техники шардирования создадим трех пользователей в системе:
 - Боб
 - Алиса
 - Генри
```shell script
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890", "name": "Bob", "surname": "Tallor", "birthday": "1994-04-10T20:21:25+00:00", "sex": "male", "city": "New Yourk", "interests": "programming"}' \
    http://localhost:10000/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890", "name": "Alice", "surname": "Swift", "birthday": "1995-10-10T20:21:25+00:00", "sex": "female", "city": "California", "interests": "running"}' \
    http://localhost:10000/auth/sign-up
curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "henry@email.com", "password": "1234567890", "name": "Henry", "surname": "Cavill", "birthday": "1993-08-19T20:21:25+00:00", "sex": "male", "city": "Washington", "interests": "sport"}' \
    http://localhost:10000/auth/sign-up
```

Получим access token'ы от системы для наших пользователей и запишем их в переменные окружения:
```shell script
export BOB_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
export ALICE_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "alice@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
export HENRY_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "henry@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
```

Проверим наличие access token-ов:
```shell script
echo $BOB_ACCESS_TOKEN
echo $ALICE_ACCESS_TOKEN
echo $HENRY_ACCESS_TOKEN
```

Теперь давайте немного початимся :). Для того, чтобы от лица Боба создать чат с Алисой, необходимо получить ее ID, так
как он понадобится для указания собеседника. Воспользуемся URL-ом на получения ID пользователя, зная его email:
```shell script
export ALICE_ID=$(curl -X GET -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    http://localhost:10000/auth/user?email=alice@email.com | jq -r '.user_id')
```

Проверим, что запрос успешно выполнился, применив команду:
```shell script
echo $ALICE_ID
```

Создадим чат от лица Боба с Алисой:
```shell script
export CHAT_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: ${BOB_ACCESS_TOKEN}" \
    -d '{"companion_id": "'"$ALICE_ID"'"}' \
    http://localhost:10000/messenger/create-chat | jq -r '.chat_id')
```

Проверим, что в переменной окружения находится UUID созданного чата:
```shell script
echo $CHAT_ID
```

Теперь отправим Алисе несколько сообщений. Для этого установим три websocket'ных соединений в трех терминальных окнах 
введем следующие команды:                               
```shell script
export BOB_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
    -d '{"email": "bob@email.com", "password": "1234567890"}' \
    http://localhost:10000/auth/sign-in | jq -r '.access_token')
websocat ws://localhost:10000/messenger/ws\?token=${BOB_ACCESS_TOKEN}
```

Теперь отправим Алисе несколько сообщений:
```shell script
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"${CHAT_ID}\", \"messages\":[{\"text\": \"Hello, Alice!\", \"status\": \"created\"}]}"}
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"${CHAT_ID}\", \"messages\":[{\"text\": \"What is up?\", \"status\": \"created\"}]}"}
{"topic":"messenger", "action": "send-message", "payload":"{\"chat_id\":\"${CHAT_ID}\", \"messages\":[{\"text\": \"I miss you!\", \"status\": \"created\"}]}"}
```

Получим со стороны Алисы, зная CHAT_ID, сообщения, которые ей отослал Боб:
```shell script
curl -X GET -H "Content-Type: application/json" -H "Authorization: ${ALICE_ACCESS_TOKEN}" \
http://localhost:10000/messenger/messages?chat_id=$CHAT_ID
```

Если все прошло успешно, должны увидеть нечто похожее: </br>
<p align="center">
    <img src="static/get_messages.png">
</p>
    
