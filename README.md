
# **简介**

devops-salt 使用c/s架构，支持服务器集群中任意节点数据采集、配置下发、批量指令执行等等动作
后续还可扩展为监控机器存活、软件告警、数据传输等等跟功能，进而可以实现dns或者软件层面的调度

# **架构**

**Agent** 

是一个使用UDP上报心跳的客户端，当同上报的心跳接收到相应之后开始下载任务，进过校验之后在机器上执, 回调执行结果

**Server** 

Server是一个监听在 UDP 端口上接收 Agent 心跳数据的网络守护进程，且运行了HTTP服务，作为任务平台，使用HTTP的方式 提交所需执行服务器的IP及指令任务

**序列化**

protobuf 作为网络序列化工具，使得传输的数据量进一步减少

# **构建**

`set GOARCH=amd64
`

`set GOOS=linux/windows`

`go build main.go`

_Server_

`docker build -t devops-salt-test:4 . `

**Linux**

目前只支持linux，稍微改动即可支持其他操作系统

# **配置 & 部署**

**Agent**

`nohup ./main -configfile /usr/local/salt/salt_client.json -server client >/dev/null 2>&1 &`

**Server**

`nohup ./main -configfile /usr/local/salt/salt_server.json -server server >/dev/null 2>&1 &
`
可直接运行docker  或者使用yaml 运行在k8s

# **联系我们**

783383650@qq.com


