#!/bin/sh

# default timezone
if [ ! -n "$TZ" ]; then
    export TZ="Asia/Shanghai"
fi

# set timezone
ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
echo $TZ > /etc/timezone 

# k8s config  switch
if [ -f "/usr/local/wechat-audio-conversion/config/.wechat-audio-conversion.yaml" ]; then
    ln -s  /usr/local/wechat-audio-conversion/config/.wechat-audio-conversion.yaml /usr/local/wechat-audio-conversion/.wechat-audio-conversion.yaml
fi

# apply config
echo "===start==="
cat /usr/local/wechat-audio-conversion/.wechat-audio-conversion.yaml
echo "====end===="

# run command
/usr/local/wechat-audio-conversion/wechat-audio-conversion api
