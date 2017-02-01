.PHONY: dev run build

dev:
	rm ./partners; go build ; \
	./partners;

run:
  ./partners

rm:
	rm bin/partners-linux-amd64; rm ~/linkit/partners-linux-amd64

build:
	export GOOS=linux; export GOARCH=amd64; \
  go build -ldflags "-s -w" -o bin/partners-linux-amd64 ; cp bin/partners-linux-amd64 ~/linkit

