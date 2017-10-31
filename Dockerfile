FROM edwinlll/wechat-audio-conversion-environment:latest 

COPY . /go/src/github.com/lvzhihao/wechat-audio-conversion 

WORKDIR /go/src/github.com/lvzhihao/wechat-audio-conversion

RUN make clean && make sbindir
RUN rm -rf /go/src/github.com/lvzhihao/wechat-audio-conversion/environment
RUN rm -f /go/src/github.com/lvzhihao/wechat-audio-conversion/.wechat-audio-conversion.yaml
RUN ln -s /usr/bin/ffmpeg sbin/ffmpeg
RUN ln -s /usr/local/sbin/silk-decoder sbin/decoder

# install go
RUN go-wrapper install

CMD ["go-wrapper", "run", "api"]
