FROM alpine:latest  

COPY . /tmp/wechat-audio-conversion-environment

RUN apk --no-cache add build-base ffmpeg && \
    cd /tmp/wechat-audio-conversion-environment/silk && make && \
    cp /tmp/wechat-audio-conversion-environment/silk/decoder /usr/local/bin/silk-decoder && \
    cp /tmp/wechat-audio-conversion-environment/silk/encoder /usr/local/bin/silk-encoder && \
    cp /tmp/wechat-audio-conversion-environment/silk/signalcompare /usr/local/bin/silk-signalcompare && \
    apk del build-base
