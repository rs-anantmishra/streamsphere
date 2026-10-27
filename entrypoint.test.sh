#!/bin/bash
# checks the URL block of entrypoint.sh (between the urls-begin/urls-end markers) for different env values
# run from the repo root: bash entrypoint.test.sh

block=$(sed -n '/# urls-begin/,/# urls-end/p' entrypoint.sh)
fail=0

# usage: check "<name>" "<protocol|applicationIp|applicationPort|contentIp|contentPort>" VAR=value ...
check() {
    local name=$1 want=$2
    shift 2
    local got
    got=$(env -i "$@" bash -c "$block"$'\n''echo "$protocol|$applicationIp|$applicationPort|$contentIp|$contentPort"' 2>&1)
    got=${got##*$'\n'}
    if [ "$got" = "$want" ]; then echo "ok   $name"; else echo "FAIL $name"; echo "  want $want"; echo "  got  $got"; fail=1; fi
}

# usage: check_says "<name>" "<text the output must contain>" VAR=value ...
check_says() {
    local name=$1 text=$2
    shift 2
    local out
    out=$(env -i "$@" bash -c "$block"$'\n''echo "$protocol|$applicationIp|$applicationPort|$contentIp|$contentPort"' 2>&1)
    if [[ "$out" == *"$text"* ]]; then echo "ok   $name"; else echo "FAIL $name"; echo "  want output containing: $text"; echo "  got: $out"; fail=1; fi
}

# usage: check_stops "<name>" "<variable the message must name>" VAR=value ...
check_stops() {
    local name=$1 var=$2
    shift 2
    local out status
    out=$(env -i "$@" bash -c "$block"$'\n''echo started' 2>&1)
    status=$?
    if [ $status -ne 0 ] && [[ "$out" == *"$var"* ]] && [[ "$out" != *started* ]]; then echo "ok   $name"; else echo "FAIL $name"; echo "  want a non-zero exit with a message naming $var, got status $status: $out"; fail=1; fi
}

# the five variables are passed to the ui as they are
check "both ports are enough" "||1282||1288" APPLICATION_PORT=1282 CONTENT_PORT=1288
check "all five are passed through" "https|ss.example.com|443|content.example.com|8443" PROTOCOL=https APPLICATION_IP=ss.example.com APPLICATION_PORT=443 CONTENT_IP=content.example.com CONTENT_PORT=8443
check "PROTOCOL accepts a trailing ://" "https||1282||1288" PROTOCOL=https:// APPLICATION_PORT=1282 CONTENT_PORT=1288

# both ports are always required
check_stops "no ports" "APPLICATION_PORT"
check_stops "no application port" "APPLICATION_PORT" CONTENT_PORT=1288
check_stops "no content port" "CONTENT_PORT" APPLICATION_PORT=1282

# HOST_IP of compose files up to 0.1.17 is read as APPLICATION_IP
check "HOST_IP is read as APPLICATION_IP" "|192.168.1.15|1282||1288" HOST_IP=192.168.1.15 APPLICATION_PORT=1282 CONTENT_PORT=1288
check "HOST_IP with CONTENT_IP keeps CONTENT_IP" "|192.168.1.15|1282|content.example.com|1288" HOST_IP=192.168.1.15 CONTENT_IP=content.example.com APPLICATION_PORT=1282 CONTENT_PORT=1288
check "APPLICATION_IP wins over HOST_IP" "|10.0.0.5|1282||1288" HOST_IP=192.168.1.15 APPLICATION_IP=10.0.0.5 APPLICATION_PORT=1282 CONTENT_PORT=1288
check_says "HOST_IP is announced as deprecated" "HOST_IP is deprecated, using it as APPLICATION_IP" HOST_IP=192.168.1.15 APPLICATION_PORT=1282 CONTENT_PORT=1288
check_says "HOST_IP next to APPLICATION_IP is announced as ignored" "HOST_IP is deprecated and ignored, APPLICATION_IP is set" HOST_IP=192.168.1.15 APPLICATION_IP=10.0.0.5 APPLICATION_PORT=1282 CONTENT_PORT=1288

# the database block (between the database-begin/database-end markers) copies the image's empty database into
# the mounted folder when that has none (#47); it runs here against a temp folder in place of /app/database
echo "--- database"
db_block=$(sed -n '/# database-begin/,/# database-end/p' entrypoint.sh)
db_root=$(mktemp -d)
trap 'chmod -R u+w "$db_root"; rm -rf "$db_root"' EXIT

# usage: run_db <case folder>; sets up <case>/template/streamsphere.db, runs the block, prints its output
run_db() {
    local root=$db_root/$1
    mkdir -p "$root/template" "$root/db"
    echo "template" > "$root/template/streamsphere.db"
    bash -c "${db_block//\/app\/database/$root}"$'\n''echo started' 2>&1
}

out=$(run_db missing)
if [ "$(cat "$db_root/missing/db/streamsphere.db" 2>/dev/null)" = "template" ] && [[ "$out" == *started* ]]; then echo "ok   a missing database is copied from the image"; else echo "FAIL a missing database is copied from the image"; echo "  got: $out"; fail=1; fi

mkdir -p "$db_root/empty/db" && : > "$db_root/empty/db/streamsphere.db"
out=$(run_db empty)
if [ "$(cat "$db_root/empty/db/streamsphere.db")" = "template" ] && [[ "$out" == *started* ]]; then echo "ok   an empty (0 byte) database is replaced"; else echo "FAIL an empty (0 byte) database is replaced"; echo "  got: $out"; fail=1; fi

mkdir -p "$db_root/kept/db" && echo "user data" > "$db_root/kept/db/streamsphere.db"
out=$(run_db kept)
if [ "$(cat "$db_root/kept/db/streamsphere.db")" = "user data" ] && [[ "$out" == *started* ]]; then echo "ok   a database with content is never touched"; else echo "FAIL a database with content is never touched"; echo "  got: $out"; fail=1; fi

mkdir -p "$db_root/readonly/db" && chmod a-w "$db_root/readonly/db"
out=$(run_db readonly)
status=$?
if [ $status -ne 0 ] && [[ "$out" != *started* ]] && [[ "$out" == *"$db_root/readonly/db"* ]] && [[ "$out" == *"cannot start"* ]]; then echo "ok   a folder it cannot write to stops the start and names the folder"; else echo "FAIL a folder it cannot write to stops the start and names the folder"; echo "  got status $status: $out"; fail=1; fi

# envsubst replaces these five names everywhere in the built bundle, not only in environment.ts, so no other
# source file may contain them in the ${name} form (a template literal would be emptied inside the container)
echo "--- bundle collisions"
clash=$(grep -rnE '\$\{(protocol|applicationIp|applicationPort|contentIp|contentPort)\}' ui/src --include=*.ts | grep -v 'environments/environment.ts' | grep -v '\.spec\.ts')
if [ -z "$clash" ]; then echo "ok   no source file collides with the envsubst names"; else echo "FAIL a source file contains an envsubst name"; echo "$clash"; fail=1; fi

exit $fail
