FROM mysql:lts

COPY *.sql /docker-entrypoint-initdb.d/

ENV LANG=C.UTF-8

EXPOSE 3306