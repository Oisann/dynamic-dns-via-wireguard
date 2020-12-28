all: build

deps:
	cat packages.txt | xargs go get

build:
	go build -v -o wg-ddns

test:
	echo "All good ;)"

clean:
	rm -f wg-ddns

.PHONY: all deps build  test clean
