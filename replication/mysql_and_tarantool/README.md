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
    - [ Причина отказа от утилиты mysql-tarantool-replication ](#task-cause)
3. [ Настройка репликации ](#replication)
    - [ Настройка master-узла MySQL ](#replication-mysql)
    - [ Настройка replicator-а ](#replication-replicator)
    - [ Настройка slave-узла tarantool ](#replication-tarantool)
4. [ Нагрузочное тестирование на чтение ](#stress-testing)
    - [ Подготовка ](#stress-testing-preparation)
    - [ Выполнение ](#stress-testing-implementation)
    - [ Результаты ](#stress-testing-results)
5. [ Итоги ](#results)

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
- [jq](https://stedolan.github.io/jq/download/) (>= version 1.5)

<a name="information-computer"></a>
### Характеристики железа
Настройка репликации, проведение нагрузочного тестирования и т.д. проводились на железе со следующими характеристиками:
- CPU - AMD Ryzen 9: 12 ядер 24 потока;
- RAM - 2xHyperX Fury Black: DDR4 DIMM 3000MHz 8GB;
- SSD - Intel® SSD 540s Series: 480GB, 2.5in SATA 6Gb/s, 16nm, TLC

<a name="task-cause"></a>
### Причина отказа от утилиты mysql-tarantool-replication
Тема довольно болезненная, однако важная для ее освещения. Постараюсь конструктивно изложить суть проблемы.
Потратив неделю безнадежного поиска причин и устранения проблем на стороне утилиты [replacator](https://github.com/tarantool/mysql-tarantool-replication),
пришел к выводу, что данная утилита полностью не актуальна и не рабочая. По следующему ряду причин:
- репозиторий попросту протухший, последний core-commit совершался более 4-х лет назад;
- были перепробованы все доступные версии MySQL на момент времени 26.11.20 и ни с одной версией [replacator](https://github.com/tarantool/mysql-tarantool-replication)
не заработал;
- [replacator](https://github.com/tarantool/mysql-tarantool-replication) собирается ТОЛЬКО на centos версии 7 и никак
иначе. В противном случае (если вы попытаетесь собрать на другой версии CentOS или на той же Ubuntu) необходимо лезть в 
сурсы и править код;
- отсутствие Docker-а (авторы прямым текстом говорят о том, что его нет и не будет - [proov](https://github.com/tarantool/mysql-tarantool-replication/pull/21));
- tutorial, который приведен [тут](https://www.tarantool.io/ru/learn/improving-mysql/), полностью противоречит тому,
что происходит в действительности (многие приведенные команды попросту не работают, а инструкции, которые приведены, как
в репозитории, так и в этой статье, сбивают с толку, потому что полностью противоположны);
- в конце концов уперся в ошибки репликтора, которые поднимались другими ребятами в issue на github'е. Проблемы не были
услышаны и maintainer'ами и никакого продвижения по их решению до сих пор нет. Авторы пометили, что это баг и казалось
бы на этом все;
- удалось переговорить с ребятами, которые остались не равнодушны к моей проблеме и откликнулись (написал на почту всем,
кто сделал fork от репы). Ответ был короткий, либо это не работает уже давным-давно, либо там надо собирать при нужной
фазе луны, поворачиваясь к северу, но смотря на юг удавалось собирать, но это было давно и сейчас никаких гарантий нет.
Эти методы я тоже пробовал также безуспешно.
- Mail в очередной раз показывает себя с ***прекрасной*** стороны и я бы все такие так же посомневался в использовании
 самого tarantool-а, не зря люди смотрят в сторону redis, думаю, по понятной причине.

Исходя из вышеописанной мной проблемы было принято решение - написать простую утилиту, которая бы позволила самым 
простым образом подписаться на event'ы из binlog'а MySQL-master-а и записать insert-инструкции в tarantool. Краткое
изложение утилиты представлено [тут](https://github.com/teploff/otus-highload/tree/main/tools/replicator).

<a name="replication"></a>
## Настройка репликации
Перед тем как перейти к настройке репликации на стороне MySQL и Tarantool необходимо поднять инфраструктуру, состоящую
из двух docker-контейнеров, а именно экземпляра MySQL и экземпляра Tarantool:
```shell script
make init
```

<a name="replication-mysql"></a>
### Настройка master-узла MySQL
Заходим в mysql-container:
```shell script
docker exec -it storage_master bash
```

Создаем папку mysql в директории /var/log/ папку mysql и даем права доступа к ней пользователю mysql:
```shell script
mkdir /var/log/mysql && chown mysql:mysql /var/log/mysql
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
bind-address = storage_master
server_id = 1
log-bin = /var/log/mysql/mysql-bin.log
binlog_format = ROW
max_binlog_size = 500M
tmpdir = /tmp
interactive_timeout=3600
wait_timeout=3600
max_allowed_packet=32M
```

Выходим из контейнера и перезапускаем его:
```shell script
exit
docker restart storage_master
```

Заходим опять в контейнер
```shell script
docker exec -it storage_master bash
```

Переходим в оболочку mysql и вводим password пароль:
```shell script
mysql -u root -p
```

Создаем пользователя для репликации и наделяем его полномочиями:
```shell script
create user 'replica'@'%' IDENTIFIED BY 'oTUSlave#2020';
GRANT REPLICATION SLAVE ON *.* TO 'replica'@'%';
```

Вызываем команду show master для того, чтобы определить MASTER_LOG_FILE и MASTER_LOG_POS, которые понадобятся нам в 
дальнейшем для запуска [replicator](https://github.com/teploff/otus-highload/tree/main/tools/replicator) утилиты:
```mysql based
show master status;
```

Результат может отличаться, но формат будет таким:</br>
<p align="center">
    <img src="static/show_master_status.png">
</p>

Выходим из оболочки MySQL
```mysql based
exit
```
и самого docker-container'а:
```shell script
exit
```

<a name="replication-replicator"></a>
### Настройка replicator-а
Перед тем, как запустить replicator, необходимо удостовериться, что переменные окружения для credentials MySQL, Tarantool
и значения BINLOG_FILE, BINLOG_POS соответствуют действительности. Для того, чтобы ознакомиться со значениями по
умолчанию, необходимо перейти в [docker-compose.yml](https://github.com/teploff/otus-highload/blob/main/replication/mysql_and_tarantool/docker-compose.yml)
и в секции **environment** сервиса **replicator** соотнести с действительностью.
Если все верно, то запускаем наш replicator следующей командой:
```shell script
make replicator
```
После чего пойдет непромедлительное реплицирование данных из MySQL в Tarantool.

Если запуск репликатора прошел успешно, то на команду:
```shell script
docker logs -f replicator_replica
```
должны увидеть следующее: </br>
<p align="center">
    <img src="static/show_replicator_launching.png">
</p>

<a name="replication-tarantool"></a>
### Настройка slave-узла tarantool
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

Выходим из docker-container:
```shell script
exit
```

Применим уже известную нам операцию накатки миграций командой:
```shell script
make migrate
```

Зайдем в tarantool и проверим, что репликация завершилась успешно и появился новый space:
```shell script
docker exec -it storage_tarantool console
box.space._space:select()
```

Должны увидеть нечто следующее:</br>
<p align="center">
    <img src="static/show_tarantool_spaces.png">
</p>

Для того, чтобы выйти из консоли, необходимо нажать **Ctrl + C**.

Выходим из container-а:
```shell script
exit
```

Если по каким-то причинам не удалось увидеть созданный space, необходимо глянуть логи replicator'а:
```shell script
docker logs replicator
```

<a name="stress-testing"></a>
## Нагрузочное тестирование на чтение
TODO

<a name="stress-testing-preparation"></a>
### Подготовка
TODO

<a name="stress-testing-implementation"></a>
### Выполнение
TODO

<a name="stress-testing-results"></a>
### Результаты
TODO

<a name="results"></a>
## Итоги
TODO
