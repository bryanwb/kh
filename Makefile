build:
	godep go build github.com/bryanwb/kh/cmd/kh
	cd fingers/hello-world && go build

clean:
	rm -f kh

test:
	godep go test

install: build
	cp -f kh /usr/local/bin/
	chmod +x /usr/local/bin/kh

all:	clean build


