#!/usr/bin/env bash

GRUNTCMD="grunt"

while getopts ":s" opt; do
  case $opt in
    s)
      GRUNTCMD="$GRUNTCMD style"
      ;;
    *)
      echo "unknown flag :("
      exit 1
  esac
done

docker compose run --rm nodejs bash -c "$GRUNTCMD"

function replace_old_asset() {
  NEW_FILE=$1
  TARGET=$2

  MD5_NEW=$(md5sum $NEW_FILE | awk '{print $1}')
  MD5_OLD=$(md5sum $TARGET | awk '{print $1}')

  echo -e "\nmd5: $MD5_NEW  \033[33m$NEW_FILE\033[0m"
  echo -e "md5: $MD5_OLD  \033[33m$TARGET\033[0m"

  if [[ "$MD5_NEW" != "$MD5_OLD" ]]; then
    cat $NEW_FILE > $TARGET
    echo -e "Replace \033[33m$TARGET\033[0m"
  fi
}

replace_old_asset public/assets/css/reprogl_temp.min.css public/assets/css/reprogl.min.css

if [[ "$GRUNTCMD" == "grunt" ]]; then
  replace_old_asset public/assets/js/reprogl_temp.min.js public/assets/js/reprogl.min.js
fi
