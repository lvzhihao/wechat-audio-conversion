OS := $(shell uname)

clean: 
	cd silk && make clean

build: */*.go
	go build

silk-decoder:
	cd silk && make
	if [ ! -d "sbin" ]; then mkdir sbin; fi;
	cp silk/decoder sbin
	chmod a+x sbin/decoder

server: silk-decoder build
	./wechat-audio-conversion api

dev: 
	DEBUG=true go run main.go api

docker-build:
	sudo docker build -t edwinlll/wechat-audio-conversion:latest .

docker-push:
	sudo docker push edwinlll/wechat-audio-conversion:latest

docker-ccr:
	sudo docker tag edwinlll/wechat-audio-conversion:latest ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion
	sudo docker push ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion
