#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

BRROOM_ID="${1:?usage: $0 <brroom_id>}"

curl_json DELETE "/brroom/$BRROOM_ID" | jq .