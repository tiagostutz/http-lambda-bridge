#!/bin/sh

http-lambda-bridge \
    --loglevel="$LOG_LEVEL"
    --proxyPass="$PROXY_PASS"
    --httpServiceInitTimeout="$HTTP_SERVICE_INIT_TIMEOUT"
    --proxyMethod="$PROXY_METHOD"