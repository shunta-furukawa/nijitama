#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
dest=build
mkdir -p "$dest/web" "$dest/curriculum"
if [ "${NIJI_VITE:-0}" = 1 ]; then
  (cd web && npm run build)
  cp -R web/dist/. "$dest/web/"
else
  cp -R web/standalone/. "$dest/web/"
  cp web/style.css "$dest/web/"
fi
cp -p Start.command "$dest/"
cp curriculum/islands.json "$dest/curriculum/"
go build -trimpath -o "$dest/nijitama" ./server
