#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
cache_root=${KOLIBRI_CACHE_DIR:-"$repo_root/.cache/kolibri"}
image_url=${KOLIBRI_IMAGE_URL:-"https://builds.kolibrios.org/en_US/data/data/kolibri.img"}
output_path="$cache_root/kolibri.img"

mkdir -p "$cache_root"

if [[ ! -f "$output_path" ]]; then
  curl -L --fail -o "$output_path" "$image_url"
fi

printf '%s\n' "$output_path"
