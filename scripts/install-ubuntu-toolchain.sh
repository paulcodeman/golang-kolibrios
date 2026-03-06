#!/usr/bin/env bash

set -euo pipefail

if ! command -v apt-get >/dev/null 2>&1; then
  echo "apt-get is required. Supported bootstrap environment: Ubuntu 24.04 or WSL Ubuntu 24.04." >&2
  exit 1
fi

sudo apt-get update
sudo apt-get install -y \
  gcc \
  gccgo \
  gcc-multilib \
  gccgo-multilib \
  make \
  nasm \
  binutils \
  mtools \
  qemu-system-x86
