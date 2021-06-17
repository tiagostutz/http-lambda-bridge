FROM tiagostutz/http-lambda-bridge:0.1.0 AS BRIDGE

FROM node:16.3-alpine3.13

# copy the executable bridge
COPY --from=BRIDGE /bin/http-lambda-bridge /bin/http-lambda-bridge

# Demo stuff
RUN npm install -g json-server
COPY demo-data.json /data/demo-data.json

# startup entrypoint script with bridge running in background
# and your http service running as originally (in foreground)
COPY startup.sh /startup.sh

RUN chmod +x /startup.sh

ENTRYPOINT ["/startup.sh"]