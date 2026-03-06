#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "$0")/.." && pwd)
app_name=templatedemo
app_dir="$repo_root/cmd/$app_name"

cleanup() {
  rm -rf "$app_dir"
}

trap cleanup EXIT

if [[ -e "$app_dir" ]]; then
  echo "temporary template check directory already exists: $app_dir" >&2
  exit 1
fi

bash "$repo_root/scripts/new-app.sh" "$app_name" "KolibriOS Template Demo"
make -C "$app_dir" clean all

if [[ ! -f "$app_dir/$app_name.kex" ]]; then
  echo "template build did not produce $app_dir/$app_name.kex" >&2
  exit 1
fi

printf 'app template check passed\n'
