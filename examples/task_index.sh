#!/usr/bin/env bash

#-- Requirements -------------------------------------------------------------------------------------------------------
command -v jq   >/dev/null 2>&2 || { echo >&2 "jq is required but is not available, aborting...";   exit 1; }
command -v curl >/dev/null 2>&2 || { echo >&2 "curl is required but is not available, aborting..."; exit 1; }

#-- Variables ----------------------------------------------------------------------------------------------------------
path=/tasks
method=GET
server=`cat scripts/api/_configuration.json | jq -r '.api_url'`

input=@./scripts/api/data/task_index.json

#-- Pre-conditions -----------------------------------------------------------------------------------------------------

#-- Action -------------------------------------------------------------------------------------------------------------
curl -X $method                                      \
     --verbose                                       \
     --data $input                                   \
     $server$path

#-- Post-Conditions ----------------------------------------------------------------------------------------------------