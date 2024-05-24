.PHONY : build clean debug

build: clean
	@echo "Build Application"
	mkdir -p build
	## Build chatbot executable
	CGO_ENABLED=0 go build -trimpath -a -ldflags '-w -s' -o ./build/chatbot ./cmd/
	## Copy default configuirations to build folder
	cp configs/configs.template.json build/configs.json
	@echo "Done"

clean:
	@echo "Clean build folder"
	rm -rf build/*

debug:
	go run ./cmd/ -c configs/debug.json 
