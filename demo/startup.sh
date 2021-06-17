#!/bin/sh

# start the lambda bridge
echo 'starting /bin/http-lambda-bridge in background'
nohup  /bin/http-lambda-bridge --logLevel="$LOG_LEVEL" --proxyPass="$PROXY_PASS" &

########################################################
## change the above line with your entrypoint command ##
########################################################
echo 'starting json-server'
json-server --watch /data/demo-data.json # start the http service
