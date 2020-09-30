FROM centos:7 AS envirenment
ENV TZ Asia/Shanghai

# 定义环境变量
ENV SCID_HOME /usr/local/salt
RUN rm -rf $SCID_HOME/
COPY ./src/main $SCID_HOME/
COPY ./etc/salt_server.json $SCID_HOME/
RUN cd /$SCID_HOME
EXPOSE 65233/udp 80
WORKDIR $SCID_HOME
CMD ["./main","-configfile", "/usr/local/salt/salt_server.json", "-server", "server"]
