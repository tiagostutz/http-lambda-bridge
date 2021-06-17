#!/bin/sh

http-lambda-bridge \
    --loglevel="$LOG_LEVEL"
    --proxyPass="$PROXY_PASS"