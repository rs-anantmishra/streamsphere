#!/bin/bash

# urls-begin
# an empty baseUrl means the UI calls the origin it was loaded from; set API_URL (e.g. https://api.example.com)
# only when the UI is hosted apart from the API
export baseUrl=${API_URL%/}
export apiUrl=$baseUrl
echo "baseUrl: ${baseUrl:-same origin}"

# content links are built by the backend (FILE_HOSTING), so they still need a host the browser can reach
# PROTOCOL is how users reach Streamsphere; content must match it or the browser blocks it as mixed content
PROTOCOL=${PROTOCOL:-http}
PROTOCOL=${PROTOCOL%%://*}://
export contentUrl=${PROTOCOL}${HOST_IP}${CONTENT_PORT:+:$CONTENT_PORT}
echo $contentUrl

export FILE_HOSTING=$contentUrl
# urls-end

# Replace placeholders in main-*.js with actual environment variable values
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