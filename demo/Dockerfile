# Get the Bridge binaries
FROM tiagostutz/http-lambda-bridge:0.1.7 AS BRIDGE

### Your original Dockerfile from here... ###
#############################################

FROM node:16.3-alpine3.13

# Demo stuff
RUN npm install -g json-server
ADD demo-data.json /data/demo-data.json

# copy the bridge executable
COPY --from=BRIDGE /bin/http-lambda-bridge /bin/http-lambda-bridge

#################################################
### ...until here, without the ENTRYPOINT or CMD ###

# ADD the startup entrypoint script having your http service running in background
# and the bridge executable running in foreground
ADD startup.sh /startup.sh

RUN chmod +x /startup.sh

ENTRYPOINT ["/startup.sh"]