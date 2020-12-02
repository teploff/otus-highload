# MySQL migration-maker tool

Which tool allows making migrations operation such as: up and down.

You should pass 3 variable via config.yaml (if you build locally) or via environment variable (if you build docker container).
Which variable are:
  - migrations_path or MIGRATIONS_PATH
  - dsn or DSN
  - operation or OPERATION

***
*migration_path* contains path directory where files with migrations are situated;
*dsn* contains destination source network - credential to connect with MySQL database;
*operation* is command for migration (up or down).