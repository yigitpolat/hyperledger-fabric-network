#!/bin/sh
mkdir -p /opt/explorer/app/platform/fabric/
mkdir -p /tmp/

mv /opt/explorer/app/platform/fabric/config.json /opt/explorer/app/platform/fabric/config.json.vanilla
cp /shared/explorer/app/config.json /opt/explorer/app/platform/fabric/config.json

cd /opt/c
node $EXPLORER_APP_PATH/main.js && tail -f /dev/null