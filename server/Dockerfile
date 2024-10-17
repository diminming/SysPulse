FROM ubuntu:24.04

RUN apt update && apt install -y libpcap-dev

RUN mkdir -p /app/syspulse/config
RUN mkdir -p /app/syspulse/bin

COPY config_staging.yaml /app/syspulse/config/server.yaml
COPY trigger.yaml /app/syspulse/config/trigger.yaml
COPY syspulse /app/syspulse/bin/syspulse_server

WORKDIR /app/syspulse

EXPOSE 24160
EXPOSE 24162

CMD ["/app/syspulse/bin/syspulse_server", "--conf", "/app/syspulse/config/server.yaml"]

