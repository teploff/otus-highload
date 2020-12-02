FROM proxysql/proxysql:2.0.13

# Update packages
RUN apt update -y --fix-missing

# Install common packages
RUN apt install -y mysql-client

# Set custom config
ADD proxysql.cnf /etc/
