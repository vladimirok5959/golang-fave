VERSION="1.0.0"

default: debug test run

debug:
	go vet ./...
	gofmt -d ./
	go build -mod vendor -o ./fave

test:
	go test ./...

run:
	@./fave -host 0.0.0.0 -port 8080 -dir ./hosts -debug true

build: clean
	@-mkdir ./bin
	@cd ./bin
	@cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.linux-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.darwin-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -mod vendor -a -o ./bin/fave.windows-amd64.exe -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	cd ./bin && find . -name 'fave*' | xargs -I{} tar czf {}.tar.gz {}
	@cp -R ./hosts/localhost/template ./bin/template
	@-rm ./bin/template/.keep
	cd ./bin && tar -zcf template.tar.gz ./template
	@-rm -r ./bin/template
	cd ./bin && shasum -a 256 * > sha256sum.txt
	cat ./bin/sha256sum.txt

clean:
	@-rm -r ./bin

format:
	gofmt -w ./

update:
	go mod vendor
	go mod download

docker-test:
	@-docker stop fave-test
	@-docker rm fave-test
	@-docker rmi fave
	docker build --rm=false --force-rm=true -t fave .
	docker run --name fave-test -d -p 8080:8080 -t -i fave:latest /app/fave.linux-amd64

docker-img:
	docker build -t fave .
