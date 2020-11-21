# Отчет о домашнем задании №6
<p align="right">
<img src="static/tarantool.png">
</p>

## Содержание
1. [ Задание ](#task)
    - [ Цель ](#task-goal)
    - [ Приобретенные навыки ](#task-skills)
    - [ Постановка задачи ](#task-statement)
2. [ Сведения ](#information)
    - [ Используемые инструменты ](#information-tools)
    - [ Характеристики железа ](#information-computer)
3. [ Настройка репликации ](#replication)
    - [ Настройка на стороне MySQL ](#replication-mysql)
    - [ Настройка на стороне Tarantool ](#replication-tarantool)

<a name="task"></a>
## Задание
Репликация из MySQL в tarantool.

<a name="task-goal"></a>
### Цель
Настроить репликацию из MySQL в tarantool и написать lua-script для осуществления выборки из tarantool-а путем 
применения хранимой процедуры.

<a name="task-skills"></a>
### Приобретенные навыки
В результате выполненного домашнего задания необходимо приобрести следующие навыки:
- администрирование MySQL;
- администрирование tarantool;
- разработка хранимых процедур для tarantool.

<a name="task-statement"></a>
### Постановка задачи
1. Выбрать любую таблицу, которую мы читаем с реплик MySQL.
2. С помощью [утилиты](https://github.com/tarantool/mysql-tarantool-replication) настроить реплицирование в tarantool
(лучше всего версии 1.10).
3. Выбрать любой запрос, переписать его на lua-процедуру и поместить его в tarantool.
4. Провести нагрузочное тестирование, сравнить tarantool и MySQL по производительности.

<a name="information"></a>
## Сведения

<a name="information-tools"></a>
### Используемые инструменты
Для выполнения дз понадобятся следующие инструменты: 
- [docker](https://docs.docker.com/get-docker/) (>= version 19.03.8) & [docker compose](https://docs.docker.com/compose/install/) (>= version 1.25.5);
- [golang](https://golang.org/doc/install) (>= version 1.15)
- [jq](https://stedolan.github.io/jq/download/) (>= version 1.5)

<a name="information-computer"></a>
### Характеристики железа
Настройка репликации, проведение нагрузочного тестирования и т.д. проводились на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC

<a name="replication"></a>
## Настройка репликации
Перед тем как перейти к настройке репликации на стороне MySQL и Tarantool необходимо поднять инфраструктуру, состоящую
из трех docker-контейнеров, а именно экземпляра MySQL, экземпляра Tarantool и экземпляра репликатора:
```shell script
make init
```

## Replicator
Заходим в replicator-container:
```shell script
docker exec -it replicator bash
```

Репликатор будет работать в виде демона systemd под названием replicatord, поэтому давайте отредактируем его служебный 
файл systemd, а именно replicatord.service, в репозитории:
```shell script
cd mysql-tarantool-replication
nano replicatord.service
```

Измените следующую строку:
```text
ExecStart=/usr/local/sbin/replicatord -c /usr/local/etc/replicatord.cfg
```

Замените расширение .cfg на .yml:
```text
ExecStart=/usr/local/sbin/replicatord -c /usr/local/etc/replicatord.yml
```

Затем скопируем некоторые файлы из репозитория replicatord в другие места locations:
```text
cp replicatord /usr/local/sbin/replicatord
cp replicatord.service /etc/systemd/system
```

### MySQL
Заходим в mysql-container:
```shell script
docker exec -it storage_mysql bash
```

Создаем папку mysql в директории */var/log/* папку mysql и даем права доступа к ней пользователю *mysql*:
```shell script
cd /var/log && mkdir mysql && chown mysql:mysql mysql
```

Устанавливаем текстовый редактор для конфигурирования, по умолчанию редактор не идет в комплектации container-а:
```shell script
apt-get update && apt-get install nano
```

Открываем конфигурацию, которая располагается по пути **/etc/mysql/conf.d/mysql.cnf**, c помощью **nano**:
```shell script
nano /etc/mysql/conf.d/mysql.cnf
```

Дописываем в секцию **[mysqld]** следующие строки:
```textmate
[mysqld]
secure-file-priv = ""
binlog_format = ROW
server_id = 1
log-bin=mysql-bin
interactive_timeout=3600
wait_timeout=3600
max_allowed_packet=32M
default_authentication_plugin=mysql_native_password
``` 

Выходим из контейнера и рестартуем его:
```shell script
docker restart storage_master
```



mkdir /var/log/tarantool
chown tarantool:tarantool /var/log/tarantool
