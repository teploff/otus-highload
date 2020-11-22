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

Выходим из container-а:
```shell script
exit 
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
exit
docker restart storage_mysql
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

Сейчас напишем стандартную Tarantool-программу путем редактирования Lua-примера, который поставляется вместе с 
Tarantool'ом:
```shell script
nano /usr/local/etc/tarantool/instances.available/example.lua
```

Полностью заменим содержимое файла следующим текстом:
``` text
box.cfg {
    listen = '*:3301';
    memtx_memory = 128 * 1024 * 1024; -- 128Mb
    memtx_min_tuple_size = 16;
    memtx_max_tuple_size = 128 * 1024 * 1024; -- 128Mb
    vinyl_memory = 128 * 1024 * 1024; -- 128Mb
    vinyl_cache = 128 * 1024 * 1024; -- 128Mb
    vinyl_max_tuple_size = 128 * 1024 * 1024; -- 128Mb
    vinyl_write_threads = 2;
    wal_mode = "none";
    wal_max_size = 256 * 1024 * 1024;
    checkpoint_interval = 60 * 60; -- one hour
    checkpoint_count = 6;
    force_recovery = true;
 
     -- 1 – SYSERROR
     -- 2 – ERROR
     -- 3 – CRITICAL
     -- 4 – WARNING
     -- 5 – INFO
     -- 6 – VERBOSE
     -- 7 – DEBUG
     log_level = 7;
     log_nonblock = true;
     too_long_threshold = 0.5;
 }
 
 local function bootstrap()
     box.schema.user.grant('guest','read,write,execute','universe')    

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

Необходимо создать символьную ссылку из instances.available (доступные экземпляры) на директорию под названием 
instances.enabled (активные экземпляры -- похоже на NGINX). В /usr/local/etc/tarantool выполните следующую команду:
```shell script
ln -s /usr/local/etc/tarantool/instances.available/example.lua instances.enabled
```

Далее мы можем запустить Lua-программу с помощью tarantoolctl (надстройки для systemd):
```shell script
tarantoolctl start example.lua
```

Переходим в console Tarantool'у, где можно проверить, что необходимые space's были успешно созданы:
```shell script
console
box.space._space:select()
```

Для того, чтобы выйти из консоли, необходимо нажать **Ctrl + D**.

### Правки
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
    user: root
    password: password
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

Далее необходимо скопировать replicatord.yml в место, где systemd будет искать его:
```shell script
cp replicatord.yml /usr/local/etc/replicatord.yml
```

Выходим из container-а и перезапускаем replicator:
```shell script
exit
docker restart replicator
```


create user 'replica'@'%' IDENTIFIED BY 'oTUSlave#2020';
GRANT REPLICATION CLIENT ON *.* TO 'replica'@'%';
GRANT REPLICATION SLAVE ON *.* TO 'replica'@'%';
GRANT SELECT ON *.* TO 'replica'@'%';
FLUSH PRIVILEGES;


sudo apt-get install cmake gcc+ libncurses5-dev libboost-dev