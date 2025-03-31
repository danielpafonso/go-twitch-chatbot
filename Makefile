.PHONY : build clean debug test

FLAGS = -trimpath -a -ldflags "-w -s"

all: clean copy build

build: app plugin
app:
	@mkdir -p build
	go build ${FLAGS} -o ./build/chatbot ./cmd/

plugin: 
	@mkdir -p build/plugins
	find plugins/commands/* -maxdepth 1 -type d -print0 | xargs -0 -I {} \
		sh -c 'd="$$(basename {})"; go build -buildmode=plugin $(FLAGS) -o ./build/plugins/"$$d".so ./{}/'
	find plugins/filters/* -maxdepth 1 -type d -print0 | xargs -0 -I {} \
		sh -c 'd="$$(basename {})"; go build -buildmode=plugin $(FLAGS) -o ./build/plugins/"$$d".so ./{}/'

clean:
	rm -rf build/*

copy:
	@mkdir -p build
	cp configs/configs.template.json build/configs.json

debug:
	go run ./cmd/ -c configs/debug.json 

test:
	go test -v ./...
