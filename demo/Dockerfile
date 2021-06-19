# Get the Bridge binaries
FROM tiagostutz/http-lambda-bridge:0.1.7 AS BRIDGE

# Your original Dockerfile
FROM node:16.3-alpine3.13

# Demo stuff
RUN npm install -g json-server
COPY demo-data.json /data/demo-data.json

# copy the bridge executable
COPY --from=BRIDGE /bin/http-lambda-bridge /bin/http-lambda-bridge

# startup entrypoint script with your http service running in background
# and the bridge running in foreground
COPY startup.sh /startup.sh

RUN chmod +x /startup.sh

ENTRYPOINT ["/startup.sh"]