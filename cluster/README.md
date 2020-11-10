0. Переходим на master
1. Т.к. никакой текстовый редактор не установлен, делаем: apt-get update && apt-get install nano
1. открываем /etc/mysql/conf.d/mysql.cnf
2. в секции [mysqld] записать следующие записи:
    - bind-address = storage_stage_1
    - server-id = 1
    - log_bin = /var/log/mysql/mysql-bin.log
3. service mysql restart
4. mysql --user user --password & create user 'replica'@'<IP>' IDENTIFIED BY 'abc@123';
5. GRANT REPLICATION SLAVE ON *.* TO 'replica'@'<IP>';
6. Show master status & save the print ;
7. Make the dump: mysqldump -u root -p --all-databases --master-data > dbdump.sql
8. Переходим на первый slave 
9. повторяем пункты 2-4
10. Загружаем dump mysql -u <user> -p < /tmp/dbdump.sql
11. stop slave;
12. CHANGE MASTER TO
    MASTER_HOST='<MASTER_IP>',
    MASTER_USER='replica',
    MASTER_PASSWORD='abc@123',
    MASTER_LOG_FILE='<FILE_NAME из принта под пунктом 6>',
    MASTER_LOG_POS='<POS из принт под пунктом 6>';
13. start slave;
14. show slave status\G ?!

     


# User generation
<img src="static/image.png" alt="">

##  What's the point of that?
This tool is based on [Faker](https://github.com/joke2k/faker) lib. It cans without troubles to generate users which consist of fields, such as:
- email;
- password;
- name;
- surname;
- birthday;
- sex;
- city;
- interests.

After generation users are stored in *.txt file for future actions.

## Using
### Requirements
- Python >= 3.6 version.
### Launching
```shell script
python main.py --count=<count of generated users> -path=<path where users should stored> 
```

Example:
```shell script
python main.py --count=100 -path=.
```
It means that there will be generate 100 users which will be sored in directory "."