#!/bin/sh

echo 'Starting background HTTP service'
##############################################################
## change the above lines with your service startup command ##
##############################################################
echo 'starting json-server'
nohup json-server --watch /data/demo-data.json  & # start the http service

# start the lambda bridge
echo 'Starting /bin/http-lambda-bridge in foreground'
/bin/http-lambda-bridge --logLevel="$LOG_LEVEL" --proxyPass="$PROXY_PASS"
