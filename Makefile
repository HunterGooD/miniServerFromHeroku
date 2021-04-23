APP?=fotocontrol
GOOS?=linux
GOARCH?=amd64

.PHONY: install build buildwebapp clean

install:
	go get ./... && \
	cd ./web && npm i

build: buildwebapp buildgo

clean: cleango cleanweb

cleango:
	@rm -rf ./dist

cleanweb:
	@rm -rf ./web/dist

buildgo: cleango
	go build -o ./dist/${APP}

buildwebapp: cleanweb
	cd ./web && \
	npm run build && \
	cd dist/ && \
	sed -i '2s/"\/web"/{{ .webAppURL }}/' ./index.html && \
	sed -i '3s/\/assets/{{ .webAppURL }}\/assets/g' ./index.html && \
	mv index.html index.tmpl && \
	rm -f index.html
