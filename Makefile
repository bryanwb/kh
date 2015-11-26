build:
	godep go build github.com/bryanwb/hand/cmd/hand

clean:
	rm -f hand

test:
	godep go test

install: build
	cp -f hand /usr/local/bin/
	chmod +x /usr/local/bin/hand

all:	clean build


