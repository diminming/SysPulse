FROM mysql:lts

COPY *.sql /docker-entrypoint-initdb.d/

ENV MYSQL_ROOT_PASSWORD=123456
ENV LANG=C.UTF-8

EXPOSE 3306