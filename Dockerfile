FROM centos:7 AS envirenment
ENV TZ Asia/Shanghai

# 定义环境变量
ENV SALT_HOME /usr/local/salt
RUN rm -rf $SALT_HOME/
COPY ./src/main $SALT_HOME/
COPY ./etc/salt_server.json $SALT_HOME/
RUN cd /$SALT_HOME
EXPOSE 65233/udp 80
WORKDIR $SALT_HOME
CMD ["./main","-configfile", "/usr/local/salt/salt_server.json", "-server", "server"]
