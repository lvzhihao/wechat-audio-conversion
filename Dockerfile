FROM edwinlll/wechat-audio-conversion-environment:latest 

WORKDIR /go/src/github.com/lvzhihao/wechat-audio-conversion

COPY . .  

# install go
RUN go-wrapper install

# clean all
RUN rm -rf *

# ext support
RUN mkdir sbin
RUN ln -s /usr/bin/ffmpeg sbin/ffmpeg
RUN ln -s /usr/local/sbin/silk-decoder sbin/decoder

CMD ["go-wrapper", "run", "api"]
