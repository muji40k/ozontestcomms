
FROM postgres:17.5

RUN mkdir -p /docker-entrypoint-initdb.d/
RUN mkdir -p /scripts/

ADD ./startup/*.sql /docker-entrypoint-initdb.d/
ADD ./*.sql /scripts/

