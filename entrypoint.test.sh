#!/bin/bash
# checks the URL block of entrypoint.sh (between the urls-begin/urls-end markers) for different env values
# run from the repo root: bash entrypoint.test.sh

block=$(sed -n '/# urls-begin/,/# urls-end/p' entrypoint.sh)
fail=0

# usage: check "<name>" "<expected baseUrl>" "<expected contentUrl>" VAR=value ...
check() {
    local name=$1 want_base=$2 want_content=$3
    shift 3
    local got
    got=$(env -i "$@" bash -c "$block"$'\n''echo "$baseUrl|$apiUrl|$contentUrl|$FILE_HOSTING"')
    got=${got##*$'\n'}
    local want="$want_base|$want_base|$want_content|$want_content"
    if [ "$got" = "$want" ]; then echo "ok   $name"; else echo "FAIL $name"; echo "  want $want"; echo "  got  $got"; fail=1; fi
}

check "no API_URL means same origin (empty)" "" "http://10.0.0.5:1288" HOST_IP=10.0.0.5 APPLICATION_PORT=1282 CONTENT_PORT=1288
check "API_URL is passed through" "https://api.example.com" "http://10.0.0.5:1288" HOST_IP=10.0.0.5 CONTENT_PORT=1288 API_URL=https://api.example.com
check "API_URL loses a trailing slash" "https://api.example.com" "http://10.0.0.5:1288" HOST_IP=10.0.0.5 CONTENT_PORT=1288 API_URL=https://api.example.com/
check "PROTOCOL=https makes the content links https" "" "https://ss.example.com:1288" HOST_IP=ss.example.com CONTENT_PORT=1288 PROTOCOL=https
check "PROTOCOL accepts a trailing ://" "" "https://ss.example.com:1288" HOST_IP=ss.example.com CONTENT_PORT=1288 PROTOCOL=https://
check "empty content port is left out" "" "https://ss.example.com" HOST_IP=ss.example.com PROTOCOL=https

exit $fail
