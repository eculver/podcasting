#!/bin/bash

set -euo pipefail

THIS_DIR="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

VOD_HOST=ftp.tdt.distributedio.netdna-cdn.com
VOD_USER=tdt.distributedio
VOD_HOME=./episodes
VOD_PASS_KEY="TDT VOD"

PZ_HOST=ftp.tdtepisodes.distributedio.netdna-cdn.com
PZ_USER=tdtepisodes.distributedio
PZ_HOME=./public_html/episodes
PZ_PASS_KEY="TDT Episodes"

function usage() {
    echo "./upload.sh <episode_mp3>"
}

function main() {
    local mp3path mp3_name vod_password pz_password

    mp3_path="$1"
    [ -z "${mp3_path}" ] && usage && exit 1
    [ ! -f "${mp3_path}" ] && echo "error: could not stat ${mp3_path}" && exit 1

    mp3_name=$(basename "${mp3_path}")
    echo "Reading $VOD_PASS_KEY keychain item"
    vod_password=$(security find-generic-password -a $LOGNAME -s "$VOD_PASS_KEY" -w)
    echo "Reading $PZ_PASS_KEY keychain item"
    pz_password=$(security find-generic-password -a $LOGNAME -s "$PZ_PASS_KEY" -w)

    pushd .
    cd "$THIS_DIR/.."

    echo "Uploading ${mp3_name} to $VOD_HOST:$VOD_HOME/${mp3_name}"
    # echo 'go run cmd/uploader/main.go --host $VOD_HOST --user $VOD_USER --pass "${vod_password}" "${mp3_path}" "$VOD_HOME/${mp3_name}"'
    go run cmd/uploader/main.go --host $VOD_HOST --user $VOD_USER --pass "${vod_password}" "${mp3_path}" "$VOD_HOME/${mp3_name}"
    echo "Uploading ${mp3_name} to $PZ_HOST:$PZ_HOME/${mp3_name}"
    # echo 'go run cmd/uploader/main.go --host $PZ_HOST --user $PZ_USER --pass "${pz_password}" "${mp3_path}" "$PZ_HOME/${mp3_name}"'
    go run cmd/uploader/main.go --host $PZ_HOST --user $PZ_USER --pass "${pz_password}" "${mp3_path}" "$PZ_HOME/${mp3_name}"

    popd
}

main "$@"
