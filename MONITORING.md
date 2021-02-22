# Мониторинг

## Содержание
1. [ Задание ](#task)
   - [ Цель ](#task-goal)
   - [ Приобретенные навыки ](#task-skills)
   - [ Постановка задачи ](#task-statement)
2. [ Сведения ](#information)
   - [ Используемые инструменты ](#information-tools)
   - [ Характеристики железа ](#information-computer)
3. [ Ход работы ](#work)
   - [ Сборка и запуск инфраструктуры ](#work-build-infrastructure)
   - [ Мониторинг с помощью Zabbix ](#work-zabbix)
        - [ Подготовка ](#work-zabbix-preparation)
        - [ Просмотр метрик ](#work-zabbix-metrics)
4. [ Итоги ](#results)

<img align="right" width="600" src="static/monitoring/preview.png">

<a name="task"></a>
## Задание
Мониторинг.

<a name="task-goal"></a>
### Цель
Организовать мониторинг сервиса диалогов.

<a name="task-skills"></a>
### Приобретенные навыки
В результате выполненного задания необходимо приобрести следующие навыки:
- эксплуатация prometheus;
- эксплуатация grafana;
- эксплуатация zabbix.

<a name="task-statement"></a>
### Постановка задачи
В процессе достижения цели необходимо:
- развернуть zabbix, prometheus и grafana;
- начать писать в prometheus бизнес-метрики сервиса чатов по принципу RED;
- начать писать в zabbix технические метрики сервера с сервисом чатов;
- организовать дашборд в grafana.

<a name="information"></a>
## Сведения
<a name="information-tools"></a>
### Используемые инструменты
Для выполнения задания понадобятся следующие инструменты:
- [docker](https://docs.docker.com/get-docker/) (>= version 19.03.8) & [docker compose](https://docs.docker.com/compose/install/) (>= version 1.25.5);
- [curl](https://curl.haxx.se/download.html) (>= version 7.68.0);

<a name="information-computer"></a>
### Характеристики железа
Задание выполнялось на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC


<a name="work"></a>
## Ход работы

<a name="work-build-infrastructure"></a>
## Сборка и запуск инфраструктуры
Клонируем наш проект:
```shell
git clone https://github.com/teploff/otus-highload.git && cd otus-highload
```

Поднимаем инфраструктуру:
```shell
make init && make migrate && make app
```


<a name="work-zabbix"></a>
## Мониторинг с помощью Zabbix

<a name="work-zabbix-preparation"></a>
### Подготовка
Перед тем как начать писать технические метрики в zabbix, необходимо знать несколько моментов.

**Инфраструктура**

Для того, чтобы слать метрики, необходим сервер zabbix-а с хранилищем для их персистентности и zabbix-frontend для 
их конфигурирования и отображения в виде dashboard-ов. Все перечисленные компоненты подняты в docker-контейнерах, а 
именно:
- в качестве хранилища был выбран [postgreSQL](https://hub.docker.com/layers/postgres/library/postgres/12-alpine/images/sha256-95ee5993459f57dddd1e42d0c11adf8363172b3828f94ed7a5ecac74da0e8ec4?context=explore) 12 версии;
- в качестве zabbix-сервера был выбран [zabbix-server-pgsql](https://hub.docker.com/r/zabbix/zabbix-server-pgsql);
- в качестве zabbix-frontend'а был выбран [zabbix-web-nginx-pgsql](https://hub.docker.com/r/zabbix/zabbix-web-nginx-pgsql).

Подробнее о том, как он собран в docker-compose'е проекта, можно посмотреть [здесь](https://github.com/teploff/otus-highload/blob/features/monitoring/deployment/docker-compose.yml#L242).


**Конфигурирование**

После того, как с инфраструктурой разобрались, еще одним не менее важным моментом для начала отправки метрик, является
само их конфигурирование на стороне zabbix-server'а. Для этого воспользуемся frontend'ом Zabbix'а. Перейдем на страницу
LogIn'а по ссылке: http://localhost:8085/. Нас попросят ввести **Username** и **Password**. В поле Username вводим
**Admin**, в поле Password - **zabbix**. Эти учетные данные являются данными по умолчанию.

Далее, находясь на вкладке **Configuration**, необходимо выбрать **Hosts**, как показано на рисунке:</br>
<p align="center">
    <img src="static/monitoring/zabbix-cfg-hosts.png">
</p>

Далее, находясь во вкладке **Hosts**, выбираем **Items**:</br>
<p align="center">
    <img src="static/monitoring/zabbix-cfg-items.png">
</p>

Затем необходимо нажать на кнопку **Create item** для создания метрики:</br>
<p align="center">
    <img src="static/monitoring/zabbix-cfg-create-button.png">
</p>

После этого, мы попадаем на форму, на которой должны создать каждую из метрик. </br>
<p align="center">
   <img src="static/monitoring/zabbix-cfg-create-form.png">
</p>

Для этого обязательно:
- выбираем в секции **Type** - **Zabbix trapper**;
- даем наименование метрики в секции **Name**;
- даем уникальный идентификатор метрики или ключ **Key**, по которому со стороны приложения будем ее отсылать;
- и тип **Type of information** самой метрики (numeric, text and ect.) 

В моем случае я создал **семь** технических метрик: четыре по **memory** и три по **cpu**. При просмотре всех метрик в
секции **Configuration** -> **Hosts** -> **Items**, пролистав метрики zabbix'а, можно увидеть только что созданные:</br>
<p align="center">
   <img src="static/monitoring/zabbix-cfg-show-items.png">
</p>

<a name="work-zabbix-metrics"></a>
### Просмотр метрик
Для того, чтобы увидеть значения метрик, необходимо перейти во вкладку **Monitoring** -> **Last data**. В моем
случае представление будет следующим:</br>
<p align="center">
    <img src="static/monitoring/zabbix-metrics-data.png">
</p>

Далее необходимо перейти в секцию **Graph**, для того, чтобы на графике увидеть динамику развития метрики. В качестве
примера выступает метрика по свободной оперативной машины **mem-free**, на которой расположен сервис диалогов: </br>
<p align="center">
    <img src="static/monitoring/zabbix-metrics-mem-free.png">
</p>

и метрика по использованию CPU **cpu-user-usage**: </br>
<p align="center">
    <img src="static/monitoring/zabbix-metrics-cpu-user.png">
</p>

Метрики собираются на стороне микросервиса диалогов с периодичностью пять секунд при помощи go-библотеки [go-osstat](https://github.com/mackerelio/go-osstat).
Сама реализация представлена [тут](https://github.com/teploff/otus-highload/blob/features/monitoring/backend/messenger/internal/infrastructure/zabbix/zabbix.go).

<a name="results"></a>
## Итоги
В ходе выполнения задания:
- был описан сбор технических метрик в zabbix;
- был реализован сбор бизнес метрик по принципу RED в prometheus;
- был организован дашборд в grafana.