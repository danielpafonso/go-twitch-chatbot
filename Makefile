.PHONY : build clean debug test

all: clean copy build

build:
	@mkdir -p build
	CGO_ENABLED=0 go build -trimpath -a -ldflags '-w -s' -o ./build/chatbot ./cmd/

clean:
	rm -rf build/*

copy:
	@mkdir -p build
	cp configs/configs.template.json build/configs.json

debug:
	go run ./cmd/ -c configs/debug.json 

test:
	go test -v ./...
