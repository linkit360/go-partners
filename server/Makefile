
rm:
	rm bin/partners-linux-amd64; rm ~/linkit/partners-linux-amd64

build:
	export GOOS=linux; export GOARCH=amd64; \
  go build -ldflags "-s -w" -o bin/partners-linux-amd64 ; cp bin/partners-linux-amd64 ~/linkit; cp dev/partners.yml ~/linkit/;

.PHONY: build