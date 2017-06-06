install:
	go get github.com/hajimehoshi/ebiten/...
	go get github.com/gopherjs/gopherjs
	go get github.com/gopherjs/webgl
	go get github.com/golang/freetype
	rm -rf node_modules
	npm install

format: # Format source code
	gofmt -w ./sample_*/

build:
	npm run build

server: # Run single instance of the server
	npm run build && cd _build && python -m SimpleHTTPServer

.PHONY: install
.PHONY: update
.PHONY: format
.PHONY: build
.PHONY: server
