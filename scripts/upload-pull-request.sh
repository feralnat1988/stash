#!/bin/sh

# assumes cross-compile.sh has already been run successfully
uploadFile() 
{
    FILE=$1
    BASENAME="$(basename "${FILE}")"
    # get available server from gofile api
    server=$(curl https://apiv2.gofile.io/getServer |cut -d "," -f 2  | cut -d "\"" -f 6)
    # abort if it takes more than two minutes to upload
    uploadedTo=$(curl -m 120 -F "email=stash@example.com" -F "file=@$FILE" "https://$server.gofile.io/uploadFile")
    resp=$(echo "$uploadedTo" | cut -d "\"" -f 4)
    URL=$(echo "$uploadedTo"|cut -d "," -f 2 | cut -d "\"" -f 6)
    if [ $resp = "ok" ] ; then
	    echo "$BASENAME uploaded to url: \"https://gofile.io/d/$URL\""
	    # print an extra newline
	    echo
    fi
}

uploadFile "dist/stash-osx"
uploadFile "dist/stash-win.exe"
uploadFile "dist/stash-linux"
