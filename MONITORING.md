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

<a name="results"></a>
## Итоги
В ходе выполнения задания:
- был описан сбор технических метрик в zabbix;
- был реализован сбор бизнес метрик по принципу RED в prometheus;
- был организован дашборд в grafana.