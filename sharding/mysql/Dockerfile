FROM mysql:8.0

# Update and clean apt cache
RUN apt-get clean && apt-get update --fix-missing && \
    # Install utils \
    apt install -y nano

# Touching log files
RUN mkdir /var/log/mysql && chown mysql:mysql /var/log/mysql
RUN touch -a /var/log/mysql/error.log && chown mysql:mysql /var/log/mysql/error.log
RUN touch -a /var/log/mysql/general.log && chown mysql:mysql /var/log/mysql/general.log
RUN touch -a /var/log/mysql/slow.log && chown mysql:mysql /var/log/mysql/slow.log

# MySQL Custom Configs
COPY mysql.cnf /etc/mysql/conf.d/mysql.cnf