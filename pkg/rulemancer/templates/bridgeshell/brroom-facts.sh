#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

BRROOM_ID="${1:?usage: $0 <brroom_id>}"

curl_json GET "/brroom/$BRROOM_ID/facts" | jq .
