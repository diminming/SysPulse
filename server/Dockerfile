FROM ubuntu:24.04

ARG ENVIRONMENT=

RUN apt-get update && apt-get install -y wget xz-utils build-essential flex bison byacc tzdata

ENV TZ=Asia/Shanghai
RUN echo "Asia/Shanghai" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata

WORKDIR /tmp
RUN wget https://www.tcpdump.org/release/libpcap-1.10.5.tar.xz && tar xvf libpcap-1.10.5.tar.xz

WORKDIR /tmp/libpcap-1.10.5
RUN ./configure && make && make install

RUN echo "/usr/local/lib" >> /etc/ld.so.conf
RUN ldconfig

RUN mkdir -p /app/syspulse/config
RUN mkdir -p /app/syspulse/bin

COPY env_conf/${ENVIRONMENT}/config.yaml /app/syspulse/config/server.yaml
COPY env_conf/${ENVIRONMENT}/trigger.yaml /app/syspulse/config/trigger.yaml
COPY syspulse /app/syspulse/bin/syspulse_server

WORKDIR /app/syspulse

EXPOSE 24160
EXPOSE 24162

CMD ["/app/syspulse/bin/syspulse_server", "--conf", "/app/syspulse/config/server.yaml"]
