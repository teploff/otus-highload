FROM haproxy:2.1.7

# Update and clean apt cache
RUN apt-get clean && apt-get update --fix-missing && \
    # Install utils \
    apt install -y default-mysql-client


# HAProxy Custom Configs
COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg