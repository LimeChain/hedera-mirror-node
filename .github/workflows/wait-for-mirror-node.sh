#!/bin/bash

mirror_node_started=false
network_identifier=""

while true;
do
    response=$(curl -sL -w "%{http_code}" -d '{"metadata":{}}' -i "http://localhost:5700/network/list")
    http_code=$(tail -n1 <<< "$response")
    if [ "$http_code" = "200" ]
    then
        echo Mirror Node has started
        network_identifier=$(tail -n2 <<< "$response" | head -n1 | jq '.network_identifiers[0]')
        mirror_node_started=true
        break
    fi

    sleep 1
done

if [[ -z $network_identifier ]]
then
    echo Rosetta Network Identifier has not been provided
    echo Exiting...
    exit 1
fi

while true;
do
    body="{ \"network_identifier\": $network_identifier, \"metadata\": {} }"
    response=$(curl -sL -w "%{http_code}" -d "$body" -i "http://localhost:5700/network/status")
    http_code=$(tail -n1 <<< "$response")

    if [ "$http_code" = "200" ]
    then
        echo Mirror Node syncing has started
        exit 0
    else
        echo Mirror Node syncing has not started yet...

        response_body=$(tail -n2 <<< "$response" | head -n1)
        is_retriable=$(echo $response_body | jq '.retriable')

        if [ "$is_retriable" = false ];
        then
            echo Request cannot be retried, response body: [$response_body]
            echo Exiting...
            exit 1
        fi
    fi

    sleep 5
done
