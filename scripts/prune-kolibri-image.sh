#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <kolibri-image-path>" >&2
  exit 1
fi

image_path=$1

# Keep the boot path and system directories intact, but strip obvious demo,
# media, and developer payload from the temporary smoke image.
prune_dirs=(
  "::/3D"
  "::/DEMOS"
  "::/DEVELOP"
  "::/FILEMA~1"
  "::/GAMES"
  "::/MEDIA"
)

prune_files=(
  "::/ALLGAMES"
  "::/DOCPACK"
  "::/EXAMPLE.ASM"
  "::/FB2READ"
  "::/HOME.PNG"
  "::/INDEX.HTM"
  "::/KUZKINA.MID"
  "::/SINE.MP3"
  "::/WELCOME.HTM"
)

for path in "${prune_dirs[@]}"; do
  mdeltree -i "$image_path" "$path" >/dev/null 2>&1 || true
done

for path in "${prune_files[@]}"; do
  mdel -i "$image_path" "$path" >/dev/null 2>&1 || true
done
