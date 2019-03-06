FROM debian

ENV FAVE_HOST=0.0.0.0 FAVE_PORT=8080 FAVE_DIR=/app/hosts FAVE_DEBUG=false

RUN apt-get update && \
 apt-get install -y wget && \
 mkdir /app && \
 mkdir /app/hosts && \
 mkdir /app/hosts/localhost && \
 mkdir /app/hosts/localhost/config && \
 mkdir /app/hosts/localhost/htdocs && \
 mkdir /app/hosts/localhost/logs && \
 mkdir /app/hosts/localhost/tmp && \
 wget -O /app/fave.linux-amd64.tar.gz https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/fave.linux-amd64.tar.gz && \
 wget -O /app/hosts/localhost/template.tar.gz https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/template.tar.gz && \
 tar -zxf /app/fave.linux-amd64.tar.gz -C /app && \
 tar -zxf /app/hosts/localhost/template.tar.gz -C /app/hosts/localhost && \
 rm /app/fave.linux-amd64.tar.gz && \
 rm /app/hosts/localhost/template.tar.gz && \
 mkdir /app/src && cp -R /app/hosts/localhost/template /app/src && \
 chmod +x /app/fave.linux-amd64

EXPOSE 8080
VOLUME /app/hosts

CMD /app/fave.linux-amd64
