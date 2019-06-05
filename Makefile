VERSION="1.0.8"

default: debug test run

debug: version dockerfile
	go vet ./...
	gofmt -d ./
	gofmt -w ./
	go build -mod vendor -o ./fave

test:
	go test ./...

run:
	@./fave -host 0.0.0.0 -port 8080 -dir ./hosts -debug true -keepalive true

build: clean version dockerfile
	@-mkdir ./bin
	@cd ./bin
	@cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.linux-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.darwin-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.windows-amd64.exe -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	cd ./bin && find . -name 'fave*' | xargs -I{} tar czf {}.tar.gz {}
	@cp -R ./hosts/localhost ./bin/localhost
	@-find ./bin/localhost -type f -name '.*' -exec rm -f {} \;
	@-find ./bin/localhost -type f -name '*.json' -exec rm -f {} \;
	@-rm ./bin/localhost/tmp/*
	cd ./bin && tar -zcf localhost.tar.gz ./localhost
	@-rm -r ./bin/localhost
	cp ./Dockerfile ./bin/Dockerfile
	cd ./bin && shasum -a 256 * > sha256sum.txt
	cat ./bin/sha256sum.txt

clean:
	@-rm -r ./bin

format:
	gofmt -w ./

update:
	go mod vendor
	go mod download

version:
	@echo "package consts" > consts/consts_version.go
	@echo "" >> consts/consts_version.go
	@echo "const ServerVersion = \"${VERSION}\"" >> consts/consts_version.go

dockerfile:
	@echo "FROM debian:latest" > Dockerfile
	@echo "" >> Dockerfile
	@echo "ENV FAVE_HOST=0.0.0.0 FAVE_PORT=8080 FAVE_DIR=/app/hosts FAVE_DEBUG=false FAVE_KEEPALIVE=true" >> Dockerfile
	@echo "" >> Dockerfile
	@echo "ADD https://github.com/vladimirok5959/golang-fave/releases/download/v${VERSION}/fave.linux-amd64.tar.gz /app/fave.linux-amd64.tar.gz" >> Dockerfile
	@echo "ADD https://github.com/vladimirok5959/golang-fave/releases/download/v${VERSION}/localhost.tar.gz /app/hosts/localhost.tar.gz" >> Dockerfile
	@echo "" >> Dockerfile
	@echo "RUN tar -zxf /app/fave.linux-amd64.tar.gz -C /app && \\" >> Dockerfile
	@echo " tar -zxf /app/hosts/localhost.tar.gz -C /app/hosts && \\" >> Dockerfile
	@echo " rm /app/fave.linux-amd64.tar.gz && \\" >> Dockerfile
	@echo " rm /app/hosts/localhost.tar.gz && \\" >> Dockerfile
	@echo " mkdir /app/src && cp -R /app/hosts/localhost /app/src && \\" >> Dockerfile
	@echo " chmod +x /app/fave.linux-amd64" >> Dockerfile
	@echo "" >> Dockerfile
	@echo "EXPOSE 8080" >> Dockerfile
	@echo "VOLUME /app/hosts" >> Dockerfile
	@echo "" >> Dockerfile
	@echo "CMD /app/fave.linux-amd64" >> Dockerfile

docker-test: dockerfile
	@-docker stop fave-test
	@-docker rm fave-test
	@-docker rmi fave
	docker build --rm=false --force-rm=true -t fave .
	docker run -d --name fave-test --cpus=".2" -m 200m -p 8080:8080 -t -i fave:latest /app/fave.linux-amd64

docker-img: dockerfile
	docker build -t fave .

docker-clr:
	@-docker stop fave-test
	@-docker rm fave-test
	@-docker rmi fave

cy-dev:
	yarn cypress open

cy:
	yarn cypress run

ab:
	ab -kc 10 -t 120 http://localhost:8080/
	ab -kc 10 -t 120 http://localhost:8080/another/
	ab -kc 10 -t 120 http://localhost:8080/not-existent-page/
	ab -kc 10 -t 120 http://localhost:8080/blog/

migrate:
	./support/migrate.sh
