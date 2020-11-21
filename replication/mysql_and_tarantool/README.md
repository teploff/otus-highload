# Отчет о домашнем задании №6
<p align="right">
<img src="static/tarantool.png">
</p>

## Содержание
1. [ Задание ](#task)
    - [ Цель ](#task-goal)
    - [ Приобретенные навыки ](#task-skills)
    - [ Постановка задачи ](#task-statement)
2. [Сведения](#information)
    - [ Используемые инструменты ](#information-tools)
    - [ Характеристики железа ](#information-computer)

<a name="information-tools"></a>
### Используемые инструменты

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