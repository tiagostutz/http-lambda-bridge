#!/bin/sh

# start the lambda bridge
nohup  /bin/http-lambda-bridge --logLevel="$LOG_LEVEL" --proxyPass="$PROXY_PASS" &


########################################################
## change the above line with your entrypoint command ##
########################################################
json-server --watch /data/demo-data.json # start the http service
