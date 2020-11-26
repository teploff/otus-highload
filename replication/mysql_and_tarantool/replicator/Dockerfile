FROM centos:7
RUN yum install -y epel-release; yum clean all

# install dependises
RUN yum install -y git ncurses-devel cmake gcc-c++ boost-devel boost boost-devel wget unzip nano bzip2 mysql-devel \
    mysql-lib expat-devel zlib-devel bzip2-devel lua-devel make /sbin/service which openssh-server
RUN git clone https://github.com/tarantool/mysql-tarantool-replication.git
RUN cd mysql-tarantool-replication && git submodule update --init --recursive && cmake . && make

#allow the ssh root access.. - Diable if you dont need but for our containers we prefer SSH access.
RUN sed -i "s/UsePAM.*/UsePAM no/g" /etc/ssh/sshd_config
RUN sed -i "s/#PermitRootLogin/PermitRootLogin/g" /etc/ssh/sshd_config

RUN echo 'root:ch@ngem3' | chpasswd

CMD ["/sbin/init"]