#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id> <fact>}"
FACT="${2:?usage: $0 <room_id> <fact>}"

curl_json GET "/room/$ROOM_ID/assert" | jq .