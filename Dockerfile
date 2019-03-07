FROM debian:latest

ENV FAVE_HOST=0.0.0.0 FAVE_PORT=8080 FAVE_DIR=/app/hosts FAVE_DEBUG=false

ADD https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/fave.linux-amd64.tar.gz /app/fave.linux-amd64.tar.gz
ADD https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/template.tar.gz /app/hosts/localhost/template.tar.gz

RUN mkdir /app/hosts/localhost/config && \
 mkdir /app/hosts/localhost/htdocs && \
 mkdir /app/hosts/localhost/logs && \
 mkdir /app/hosts/localhost/tmp && \
 tar -zxf /app/fave.linux-amd64.tar.gz -C /app && \
 tar -zxf /app/hosts/localhost/template.tar.gz -C /app/hosts/localhost && \
 rm /app/fave.linux-amd64.tar.gz && \
 rm /app/hosts/localhost/template.tar.gz && \
 mkdir /app/src && cp -R /app/hosts/localhost/template /app/src && \
 chmod +x /app/fave.linux-amd64

EXPOSE 8080
VOLUME /app/hosts

CMD /app/fave.linux-amd64
