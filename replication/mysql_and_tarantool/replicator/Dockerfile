FROM centos:7
RUN yum install -y epel-release; yum clean all

# install dependises
RUN yum install -y git ncurses-devel cmake gcc-c++ boost-devel boost boost-devel wget unzip nano bzip2 mysql-devel \
    mysql-lib expat-devel zlib-devel bzip2-devel lua-devel make
RUN git clone https://github.com/tarantool/mysql-tarantool-replication.git
RUN cd mysql-tarantool-replication && git submodule update --init --recursive && cmake . && make

CMD ["/bin/sh", "-c", "while true; do echo ping >> /dev/null; sleep 60; done"]