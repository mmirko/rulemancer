#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

BRIDGE_ID="${1:?usage: $0 <bridge_id|bridge_name>}"

curl_json GET "/bridge/$BRIDGE_ID" | jq .
