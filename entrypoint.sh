#!/bin/bash

APP_PROTOCOL='http://'
export apiUrl=${APP_PROTOCOL}${HOST_IP}:${APPLICATION_PORT}
echo $apiUrl
export baseUrl=${APP_PROTOCOL}${HOST_IP}:${APPLICATION_PORT}
echo $baseUrl
export contentUrl=${APP_PROTOCOL}${HOST_IP}:${CONTENT_PORT}
echo $contentUrl

export FILE_HOSTING=$contentUrl

# Replace placeholders in main-*.js with actual environment variable values
if [ -z "$apiUrl" ]; then
  echo "apiUrl environment variable is not set. Application cannot start."
  exit 1;
fi

if [ -z "$baseUrl" ]; then
  echo "baseUrl environment variable is not set. Application cannot start."
  exit 1;
fi

if [ -z "$contentUrl" ]; then
  echo "contentUrl environment variable is not set. Application cannot start."
  exit 1;
fi

cd /app/browser
mainfile=$(ls main-*.js)
# echo "mainfile is: $mainfile"
# Run envsubst to replace placeholders in main-*.jx with actual environment variable values
envsubst '${apiUrl},${baseUrl},${contentUrl}' < ./${mainfile} > ./${mainfile}.tmp && \
mv ./${mainfile}.tmp ./${mainfile}

# update yt-dlp
cd /app/utils
./yt-dlp_linux -U

# run streamsphere
cd /app/application
chmod +x streamsphere
./streamsphere