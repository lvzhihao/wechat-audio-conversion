OS := $(shell uname)

all: bindir ffmpeg silk-decoder build
	go test -v

bindir:
	if [ ! -d "bin" ]; then mkdir bin; fi;

# for ubuntu only~~~ producer use docker image
ffmpeg: bindir
	if [ ! -f "bin/ffmpeg" ]; then \
	 apt-get install ffmpeg -y && ln -s /usr/bin/ffmpeg bin/ffmpeg; \
	fi;

silk-decoder: bindir
	if [ ! -f "bin/decoder" ]; then \
	cd environment/silk && make && cd ../../ && cp environment/silk/decoder bin && chmod a+x bin/decoder; \
	fi;

clean: 
	cd environment/silk && make clean
	rm -rf bin
	rm -f wechat-audio-conversion

build: */*.go
	go build

server: all
	./wechat-audio-conversion api

dev: bindir ffmpeg silk-decoder
	DEBUG=true go run main.go api

env-build:
	cd environment &&  docker build -t edwinlll/wechat-audio-conversion-environment:latest .

env-push:
	 docker push edwinlll/wechat-audio-conversion-environment:latest

docker-build:
	 docker build -t edwinlll/wechat-audio-conversion:latest .

docker-push:
	 docker push edwinlll/wechat-audio-conversion:latest

docker-ccr:
	 docker tag edwinlll/wechat-audio-conversion:latest ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion:latest
	 docker push ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion:latest
	 docker rmi ccr.ccs.tencentyun.com/wdwd/wechat-audio-conversion:latest

docker-uhub:
	 docker tag edwinlll/wechat-audio-conversion:latest uhub.service.ucloud.cn/mmzs/wechat-audio-conversion:latest
	 docker push uhub.service.ucloud.cn/mmzs/wechat-audio-conversion:latest
	 docker rmi uhub.service.ucloud.cn/mmzs/wechat-audio-conversion:latest

docker-ali:
	 docker tag edwinlll/wechat-audio-conversion:latest registry.cn-hangzhou.aliyuncs.com/weishangye/wechat-audio-conversion:latest
	 docker push registry.cn-hangzhou.aliyuncs.com/weishangye/wechat-audio-conversion:latest
	 docker rmi registry.cn-hangzhou.aliyuncs.com/weishangye/wechat-audio-conversion:latest
