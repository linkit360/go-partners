.PHONY: run build

VERSION=$(shell git describe --always --long --dirty)

version:
	 @echo Version IS $(VERSION)

rm:
	rm bin/partners-linux-amd64; rm ~/linkit/partners-linux-amd64

build:
	export GOOS=linux; export GOARCH=amd64; \
	sed -i "s/%VERSION%/$(VERSION)/g" /home/centos/vostrok/utils/metrics/metrics.go; \
  go build -ldflags "-s -w" -o bin/partners-linux-amd64 ; cp bin/partners-linux-amd64 ~/linkit; cp dev/partners.yml ~/linkit/;

.PHONY: build