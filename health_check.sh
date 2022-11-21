#!/bin/bash
httport=$HTTPPORT

if [ -z "$httport" ]; then
    httport=80
fi

url="http://localhost:$httport/healthz"
code=$(curl -sIL -w "%{http_code}\n" -o /dev/null $url)
echo $code
# -eq 200 检测域名是否返回200
if [ $code -eq 200 ]; then
    echo "OK"
    exit 1
else
    echo "ERROR"
    exit 0
fi