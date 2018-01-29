FROM  centos:centos7.1.1503

ENV TZ "Asia/Shanghai"

COPY  aliyun-mirror.repo /etc/yum.repos.d/CentOS-Base.repo
COPY  aliyun-epel.repo /etc/yum.repos.d/epel.repo
COPY  mogodb.repo /etc/yum.repos.d
# RUN sed -i 's/plugins=1/plugins=0/' /etc/yum.conf 
# RUN sed -i 's/enabled=1/enabled=0/' /etc/yum/pluginconf.d/fastestmirror.conf
RUN rpm --rebuilddb && yum install -y deltarpm
# 配置go环境
RUN rpm --rebuilddb &&  yum install -y go git  wget tar  vim-enhanced  supervisor 


RUN go get -v github.com/gin-gonic/gin
RUN go get -v github.com/garyburd/redigo/redis
RUN go get -v gopkg.in/mgo.v2
RUN go get -v github.com/Shopify/sarama
RUN go get -v github.com/gorilla/websocket
# 配置supervisor

COPY  supervisord.conf /etc/supervisord.conf


RUN mkdir -p /etc/supervisor.conf.d && \
    mkdir -p /var/log/supervisor 

# 安装redis mongodb
RUN rpm --rebuilddb && yum install -y mongodb-org && \
    rpm --rebuilddb && yum install -y redis 


COPY  supervisor_go.conf /etc/supervisor.conf.d/kafka.conf

#配置jdk
RUN rpm --rebuilddb &&  yum install -y java-1.8.0-openjdk.x86_64 && \
    yum clean all

#配置kafka
RUN  wget https://mirrors.cnnic.cn/apache/kafka/1.0.0/kafka_2.11-1.0.0.tgz && \
     tar xf kafka_2.11-1.0.0.tgz -C /usr/local         && \
     rm -rf ./kafka_2.11-1.0.0.tgz  && \
     ln -sv /usr/local/kafka_2.11-1.0.0 /usr/local/kafka 
  
EXPOSE 8000

ENTRYPOINT ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
