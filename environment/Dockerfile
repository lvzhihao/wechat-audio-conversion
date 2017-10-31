FROM golang:1.9

COPY . /tmp/wechat-audio-conversion-environment

# Install ffmpeg 
RUN echo 'deb http://www.deb-multimedia.org stretch main non-free' >> /etc/apt/sources.list && \
        apt-get update && \
        apt-get install -y --force-yes \
        deb-multimedia-keyring \
        ffmpeg \
        gettext-base  && \
        apt-get autoremove -yqq && \
        apt-get clean -yqq

# Install silk-decoder
RUN cd /tmp/wechat-audio-conversion-environment/silk && make
RUN cp /tmp/wechat-audio-conversion-environment/silk/decoder /usr/local/sbin/silk-decoder

# Clean all
RUN rm -rf /tmp/wechat-audio-conversion-environment

WORKDIR $GOPATH
