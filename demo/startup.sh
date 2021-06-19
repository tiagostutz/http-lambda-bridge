#!/bin/sh

# start the lambda bridge in background
echo 'Starting /bin/http-lambda-bridge in background'
nohup /bin/http-lambda-bridge --logLevel="$LOG_LEVEL" --proxyPass="$PROXY_PASS" --httpServiceInitTimeout=$HTTP_SERVICE_INIT_TIMEOUT  & # start the http service

echo 'Starting HTTP service as the main process'

##############################################################
## change the above lines with your service startup command ##
##############################################################
echo 'starting json-server in foreground'
json-server --watch /data/demo-data.json
