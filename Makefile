OS := $(shell uname)

all: sbindir ffmpeg silk-decoder build
	go test -v

sbindir:
	if [ ! -d "sbin" ]; then mkdir sbin; fi;

# for ubuntu only~~~ producer use docker image
ffmpeg: sbindir
	if [ ! -f "sbin/ffmpeg" ]; then \
	sudo apt-get install ffmpeg -y && ln -s /usr/bin/ffmpeg sbin/ffmpeg; \
	fi;

silk-decoder: sbindir
	if [ ! -f "sbin/decoder" ]; then \
	cd environment/silk && make && cd ../../ && cp environment/silk/decoder sbin && chmod a+x sbin/decoder; \
	fi;

clean: 
	cd environment/silk && make clean
	rm -rf sbin
	rm -f wechat-audio-conversion

build: */*.go
	go build

server: all
	./wechat-audio-conversion api

dev: sbindir ffmpeg silk-decoder
	DEBUG=true go run main.go api

env-build:
	cd environment && sudo docker build -t edwinlll/wechat-audio-conversion-environment:latest .

env-push:
	sudo docker push edwinlll/wechat-audio-conversion-environment:latest

docker-build:
	sudo docker build -t edwinlll/wechat-audio-conversion:latest .

docker-push:
	sudo docker push edwinlll/wechat-audio-conversion:latest

docker-ccr:
	sudo docker tag edwinlll/wechat-audio-conversion-environment:latest ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion-environment
	sudo docker push ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion-environment
	sudo docker tag edwinlll/wechat-audio-conversion:latest ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion
	sudo docker push ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion
