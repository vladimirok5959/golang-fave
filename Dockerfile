FROM debian:latest

ENV FAVE_HOST=0.0.0.0 FAVE_PORT=8080 FAVE_DIR=/app/hosts FAVE_DEBUG=false

ADD https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/fave.linux-amd64.tar.gz /app/fave.linux-amd64.tar.gz
ADD https://github.com/vladimirok5959/golang-fave/releases/download/v1.0.0/localhost.tar.gz /app/hosts/localhost.tar.gz

RUN tar -zxf /app/fave.linux-amd64.tar.gz -C /app && \
 tar -zxf /app/hosts/localhost.tar.gz -C /app/hosts && \
 rm /app/fave.linux-amd64.tar.gz && \
 rm /app/hosts/localhost.tar.gz && \
 mkdir /app/src && cp -R /app/hosts/localhost /app/src && \
 chmod +x /app/fave.linux-amd64

EXPOSE 8080
VOLUME /app/hosts

CMD /app/fave.linux-amd64
