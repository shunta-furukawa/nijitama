#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
cp web/style.css web/standalone/style.css
export NIJI_WEB=web/standalone
exec go run ./server
