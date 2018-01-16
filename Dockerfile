FROM golang:1.9 as builder
WORKDIR /go/src/github.com/lvzhihao/wechat-audio-conversion
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM edwinlll/wechat-audio-conversion-environment:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /usr/local/wechat-audio-conversion
COPY --from=builder /go/src/github.com/lvzhihao/wechat-audio-conversion/wechat-audio-conversion .
COPY .wechat-audio-conversion.yaml.sample .wechat-audio-conversion.yaml
# ext support
RUN mkdir bin && \
    ln -s /usr/bin/ffmpeg bin/ffmpeg && \
    ln -s /usr/local/bin/silk-decoder bin/decoder
ENV PATH /usr/local/wechat-audio-conversion:$PATH
CMD ["wechat-audio-conversion", "api"]
