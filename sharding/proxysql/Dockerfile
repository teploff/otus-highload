FROM proxysql/proxysql:latest

# Update packages
RUN apt update -y --fix-missing

# Install common packages
RUN apt-get update & apt-get install -y mysql-client nano

# Set custom config
ADD proxysql.cnf /etc/
