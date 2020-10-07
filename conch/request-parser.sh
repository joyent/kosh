#!/bin/sh
curl -o $1.json https://conch.joyent.us/schema/request/$1
schematyper -o types/RequestType_${1}.go --package="types" --ptr-for-omit $1.json
rm $1.json
