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

Инвалидация по событию
Удаляем данные из кеша при изменении в базе

docker exec -it redis-cache redis-cli
AUTH password
FLUSHALL

```shell
cd backend/tools/cache-heater-enabler
go run main.go --addr="localhost:4222" --cluster_id="stan-cluster" --subject="cache-reload"
```

<a name="results"></a>
## Итоги
В ходе выполнения домашнего задания: