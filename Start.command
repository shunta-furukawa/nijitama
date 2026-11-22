#!/bin/bash
# Finder and Terminal launcher for macOS.
set -euo pipefail
umask 077
cd "$(dirname "$0")"
export NIJI_DATA="${NIJI_DATA:-$HOME/Library/Application Support/nijitama}"
mkdir -p "$NIJI_DATA"
pinfile="$NIJI_DATA/parent.pin"
if [[ -z ${NIJI_PARENT_PIN:-} ]]; then
  if [[ ! -f $pinfile ]]; then
    echo '保護者番号を設定します。既存の学習データがある場合は、以前と同じ番号を入力してください。'
    read -rsp '保護者番号（数字6桁以上）: ' pin; echo
    read -rsp 'もう一度: ' confirm; echo
    if [[ ! $pin =~ ^[0-9]{6,}$ || $pin != "$confirm" ]]; then
      echo '番号が一致しないか、形式が正しくありません。' >&2
      exit 1
    fi
    tmp=$(mktemp "$NIJI_DATA/.parent-pin.XXXXXX")
    trap 'rm -f "$tmp"' EXIT
    printf '%s\n' "$pin" > "$tmp"
    ln "$tmp" "$pinfile" || exit 1
    rm -f "$tmp"
    trap - EXIT
  fi
  NIJI_PARENT_PIN=$(cat "$pinfile")
fi
export NIJI_PARENT_PIN
if [[ -z ${NIJI_LAN_IP:-} ]]; then
  interface=$(networksetup -listallhardwareports | awk '/^Hardware Port: (Wi-Fi|AirPort)$/{wifi=1;next} wifi && /^Device:/{print $2; exit}')
  if [[ -z $interface ]] || ! NIJI_LAN_IP=$(ipconfig getifaddr "$interface"); then
    echo 'Wi-FiのIPv4を取得できません。Wi-Fiに接続するか、NIJI_LAN_IP に家庭内IPv4を指定してください。' >&2
    exit 1
  fi
fi
export NIJI_LAN_IP
export NIJI_LOCAL_NAME="$(scutil --get LocalHostName 2>/dev/null || true)"
if [[ -f go.mod ]]; then
  ./scripts/build.sh
  cd build
fi
exec caffeinate -i ./nijitama
