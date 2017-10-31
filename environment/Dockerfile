FROM golang:1.9

COPY . /tmp/wechat-audio-conversion-environment

WORKDIR /tmp/wechat-audio-conversion-environment

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
RUN cd silk && make
RUN cp silk/decoder /usr/local/sbin/silk-decoder
RUN cd silk && make clean

# Clean all
RUN cd .. && rm -rf ./wechat-audio-conversion-environment
