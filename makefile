all: phishflood

phishflood: build
	@go build -o build/phishflood cmd/phishflood/*.go

install: phishflood
	@cp build/phishflood /usr/bin/
	@chmod a+x /usr/bin/phishflood

build:
	@mkdir -p build

clean:
	@rm -rf build
