#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id> <relation>}"
RELATION="${2:?usage: $0 <room_id> <relation>}"

curl_json POST "/room/$ROOM_ID/query/$RELATION" | jq .