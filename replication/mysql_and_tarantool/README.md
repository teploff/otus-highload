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

В большинстве своем случае настройка, будет происходить согласно [статье](https://www.tarantool.io/ru/learn/improving-mysql/), 
однако статья не рассчитана на docker-container'ы и в ней есть ряд моментов, которые не работают.
Ниже приводится подробное описание настройки каждого из container-ов:

### MySQL
Заходим в mysql-container:
```shell script
docker exec -it storage_mysql bash
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
bind-address = storage_mysql
server_id = 1
log-bin=mysql-bin
binlog_format = ROW
interactive_timeout=3600
wait_timeout=3600
max_allowed_packet=32M
```


Выходим из контейнера и рестартуем его:
```shell script
exit
docker restart storage_mysql
```

```shell script
docker exec -it storage_mysql bash
```

```shell script
mysql -u root -p
```

```shell script
create user 'replica'@'%' IDENTIFIED BY 'oTUSlave#2020';
GRANT REPLICATION CLIENT ON *.* TO 'replica'@'%';
GRANT REPLICATION SLAVE ON *.* TO 'replica'@'%';
GRANT SELECT ON *.* TO 'replica'@'%';
FLUSH PRIVILEGES;
```

```mysql based
show master status;
```

```mysql based
CREATE DATABASE menagerie;
use menagerie;
DROP TABLE IF EXISTS pet;

CREATE TABLE pet (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name2   VARCHAR(20),
  owner   VARCHAR(20),
  species VARCHAR(20)
);

```

Выходим из container-а:
```shell script
exit
exit
```

### Replicator
Заходим в replicator-container:
```shell script
docker exec -it replicator bash
```

Репликатор будет работать в виде демона systemd под названием replicatord, поэтому давайте отредактируем его служебный 
файл systemd, а именно replicatord.service, в репозитории:
```shell script
cd mysql-tarantool-replication && nano replicatord.service
```

Измените следующую строку:
```text
...
ExecStart=/usr/local/sbin/replicatord -c /usr/local/etc/replicatord.cfg
...
```

Замените расширение .cfg на .yml:
```text
...
ExecStart=/usr/local/sbin/replicatord -c /usr/local/etc/replicatord.yml
...
```

Затем скопируем некоторые файлы из репозитория replicatord в другие места locations:
```text
cp replicatord /usr/local/sbin/replicatord
cp replicatord.service /etc/systemd/system
```

Выходим из container-а:
```shell script
exit 
```

### Tarantool
Заходим в tarantool-container:
```shell script
docker exec -it storage_tarantool sh
```

Создаем папку tarantool в директории */var/log/* и даем права доступа к ней пользователю *tarantool*:
```shell script
mkdir /var/log/tarantool && chown tarantool:tarantool /var/log/tarantool
```

Устанавливаем текстовый редактор для конфигурирования, по умолчанию редактор не идет в комплектации container-а:
```shell script
apk update && apk add nano
```

Переходим в console:
```shell script
console
```

Выполним следующие команды:
``` shell script
 box.schema.user.grant('guest','read,write,execute','universe')

 function bootstrap()   
     if not box.space.mysqldaemon then
         s = box.schema.space.create('mysqldaemon')
         s:create_index('primary',
         {type = 'tree', parts = {1, 'unsigned'}})
     end
 
     if not box.space.mysqldata then
         t = box.schema.space.create('mysqldata')
         t:create_index('primary',
         {type = 'tree', parts = {1, 'unsigned'}})
     end
 
 end
 
 bootstrap()
```

Проверим, создались ли два space-а(mysqldaemon и mysqldata):
```shell script
box.space._space:select()
```

Для того, чтобы выйти из консоли, необходимо нажать **Ctrl + C**.

Выходим из container-а:
```shell script
exit
```

### Следующие шаги
Заходим в replicator-container:
```shell script
docker exec -it replicator bash
```

Переходим к конфигурации репликатора:
```shell script
cd mysql-tarantool-replication
nano replicatord.yml
```

Заменяем ее полностью на:
```yaml
mysql:
    host: storage_mysql
    port: 3306
    user: replica
    password: oTUSlave#2020
    connect_retry: 15 # seconds

tarantool:
    host: storage_tarantool:3301
    binlog_pos_space: 512
    binlog_pos_key: 0
    connect_retry: 15 # seconds
    sync_retry: 1000 # milliseconds

mappings:
    - database: menagerie
      table: pet
      columns: [ id, name2, owner, species ]
      space: 513
      key_fields: [ 0 ]

spaces:
    513:
        replace_null:
            1: { string: "" }
            2: { string: "" }
```

Далее необходимо скопировать replicatord.yml туда, где systemd будет искать его:
```shell script
cp replicatord.yml /usr/local/etc/replicatord.yml
systemctl start replicatord
```

Смотрим логги:
```shell script
tail -f /var/log/replicatord.log
```

В отдельном terminal-е заходим в MySQL-container и сделаем запись:
```mysql based
docker exec -it storage_mysql bash
mysql -u root -p
use menagerie;
INSERT INTO pet(name2, owner, species) VALUES ('Spot', 'Brad', 'dog');
```

В отдельном terminal-е заходим в Tarantool-console и проверим, среплицировалась ли запись:
```shell script
docker exec -it storage_tarantool console
box.space._space:select()
```