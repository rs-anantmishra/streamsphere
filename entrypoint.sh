#!/bin/bash

# urls-begin
# HOST_IP of compose files up to 0.1.17 is read as APPLICATION_IP
if [ -n "$HOST_IP" ] && [ -z "$APPLICATION_IP" ]; then
  echo "HOST_IP is deprecated, using it as APPLICATION_IP"
  APPLICATION_IP=$HOST_IP
elif [ -n "$HOST_IP" ]; then
  echo "HOST_IP is deprecated and ignored, APPLICATION_IP is set"
fi

# both ports are always required: the left sides of the two "ports" mappings
if [ -z "$APPLICATION_PORT" ]; then
  echo "APPLICATION_PORT is not set. Set it to the port you open Streamsphere on (the left side of the \"ports\" mapping). Application cannot start."
  exit 1
fi
if [ -z "$CONTENT_PORT" ]; then
  echo "CONTENT_PORT is not set. Set it to the port of the content server (the left side of the content \"ports\" mapping). Application cannot start."
  exit 1
fi

# what is not set (protocol, the ips) is taken from the address the page was opened at, by the ui
export protocol=${PROTOCOL%%://*}
export applicationIp=$APPLICATION_IP
export applicationPort=$APPLICATION_PORT
export contentIp=$CONTENT_IP
export contentPort=$CONTENT_PORT
echo "protocol: ${protocol:-same as the page}"
echo "api: ${applicationIp:-page host}:${applicationPort}"
echo "content: ${contentIp:-${applicationIp:-page host}}:${contentPort}"
# urls-end

cd /app/browser
mainfile=$(ls main-*.js)
# echo "mainfile is: $mainfile"
# Run envsubst to replace placeholders in main-*.jx with actual environment variable values
envsubst '${protocol},${applicationIp},${applicationPort},${contentIp},${contentPort}' < ./${mainfile} > ./${mainfile}.tmp && \
mv ./${mainfile}.tmp ./${mainfile}

# update yt-dlp
cd /app/utils
./yt-dlp_linux -U

# run streamsphere
cd /app/application
chmod +x streamsphere
./streamsphere