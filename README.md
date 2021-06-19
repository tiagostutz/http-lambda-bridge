# HTTP Lambda Bridge

Easily convert any HTTP Rest service into a Lambda application

If you plan to migrate your monolith to AWS Lambda Functions, but can't wait until all the functions are ready to be deployed to have the benefits of the Lambda platform, this bridge is for you!

## Getting Started

You can check the [demo](demo) folder for a complete example on how to host a [json-server](https://www.npmjs.com/package/json-server) as a lambda function.

## Serving your Docker Restful API monolith as an AWS Lambda 

First, create an `startup.sh` file that starts the lambda bridge in background and your HTTP Rest service as the main process of the image. Like this:
```bash
#!/bin/sh

# start the lambda bridge in background
echo 'Starting /bin/http-lambda-bridge in background'
nohup /bin/http-lambda-bridge --logLevel="$LOG_LEVEL" --proxyPass="$PROXY_PASS" --httpServiceInitTimeout=$HTTP_SERVICE_INIT_TIMEOUT  & # start the http service

echo 'Starting HTTP service as the main process'

##############################################################
## change the above lines with your service startup command ##
##############################################################
echo 'starting HTTP Rest service in foreground'
# YOUR STARTUP COMMANDO HERE, LIKE `npm start` or `./my-go-service`
```

Second, build your Dockerfile like this:

```Dockerfile
FROM tiagostutz/http-lambda-bridge:0.1.7 AS BRIDGE

### Your original Dockerfile from here... ###
#############################################

# ...

#################################################
### ...until here, without the ENTRYPOINT or CMD ###

# copy the bridge executable
COPY --from=BRIDGE /bin/http-lambda-bridge /bin/http-lambda-bridge

# ADD the startup entrypoint script having your http service running in background
# and the bridge executable running in foreground
ADD startup.sh /startup.sh

RUN chmod +x /startup.sh

ENTRYPOINT ["/startup.sh"]

```

After building and deploying your function (see [build-deploy-aws-lambda.sh](demo/build-deploy-aws-lambda.sh) in the demo) go ahead and `curl` your HTTP function address to see.

### Environment Variables

- LOG_LEVEL: trace | debug | info | warning | error | panic
- PROXY_PASS: local endpoint of the HTTP Rest service bridging to
- HTTP_SERVICE_INIT_TIMEOUT: maximum timeout to wait for the HTTP Rest service to be up and running


## Important note!

In order for the bridge to receive the complete request path from the API Gateway, you need to override the PayloadFormatVersion of the APIGateway like this:

```yaml

  HttpAPIGatewayOverrides:
    Type: AWS::ApiGatewayV2::ApiGatewayManagedOverrides
    Properties: 
      ApiId: !Ref HttpAPIGateway
      Integration: 
        PayloadFormatVersion: 1.0
    ...
```

Check the [cf-api-lambda.yml](demo/cf-api-lambda.yml) for more details

