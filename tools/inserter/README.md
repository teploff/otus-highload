# Insert generated user in MySQL using python
<img align="right" width="250" src="static/image.png" alt="">

##  What's the point of that?
This tool is based on [PyMySQL](https://github.com/PyMySQL/PyMySQL) lib for inserting in MySQL database.
Inserting is based on generated users which are given by [generator tool](https://github.com/teploff/otus-highload/tree/main/tools/generator).

## Using
### Requirements
- Python >= 3.6 version;
- generated users by [tool] (https://github.com/teploff/otus-highload/tree/main/tools/generator);
- MySQL server.

### Launching
```shell script
python main.py -cfg=<path to your config> -path=<path where users should stored in *.txt files>  -size=<size of butch for inserting>
```

Example:
```shell script
python main.py -cfg=./config.yaml -path=../snapshot -size=10000 
```
It means that there will be *10k* inserts of users simultaneously (size of butch) which are stored in *../snapshot* directory with a configuration which path is  */config.yaml*

Configuration example:
storage:
  host: ${env:STORAGE_HOST, localhost}
  port: ${env:STORAGE_PORT, 3306}
  user: ${env:STORAGE_USER, user}
  password: ${env:STORAGE_PASSWORD, password}
  db: ${env:STORAGE_DB, social-network}
  charset: ${env:STORAGE_CHARSET, utf8mb4}
  
It consists main credentials for database connection. First of all configuration will try find out credentials in environment variables, such as:
- STORAGE_HOST;
- STORAGE_PORT;
- STORAGE_USER;
- STORAGE_PASSWORD;
- STORAGE_DB;
- STORAGE_CHARSET.

But if they will be absent configuration will take their from default values which are given on another part of configuration, such as:
- localhost;
- 3306;
- user;
- password;
- social-network;
- utf8mb4.

## NOTE
In *../generator/snapshot* directory there are 1m compressed generated users in *.tar.gz archives. If you want use it - before launch inserting you should uncompressed their.