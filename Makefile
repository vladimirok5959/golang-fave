VERSION="1.0.1"

default: debug

debug:
	@go vet ./...
	@go build -o ./fave

build: clean
	@-mkdir ./bin
	@cd ./bin
	@cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/fave.linux-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ./bin/fave.darwin-amd64 -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ./bin/fave.windows-amd64.exe -ldflags='-X main.Version=$(VERSION) -extldflags "-static"'
	cd ./bin && find . -name 'fave*' | xargs -I{} tar czf {}.tar.gz {}
	cd ./bin && shasum -a 256 * > sha256sum.txt
	cat ./bin/sha256sum.txt

clean:
	@-rm -r ./bin

test:
	@go test ./...

run:
	@./fave -host 0.0.0.0 -port 8080 -dir ./hosts
