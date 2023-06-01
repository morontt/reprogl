#!/usr/bin/env bash

GRUNTCMD="grunt"
CURRENTDATE=$(date '+%y%j%H%M%S')

sed -i 's/reprogl\.min\.v\([0-9]\+\)\.css/reprogl\.min\.v'$CURRENTDATE'\.css/' ./templates/layout/base.gohtml

while getopts ":s" opt; do
  case $opt in
    s)
      GRUNTCMD="$GRUNTCMD style"
      ;;
  esac
done

if [[ "$GRUNTCMD" == "grunt" ]]; then
  sed -i 's/reprogl\.min\.v\([0-9]\+\)\.js/reprogl\.min\.v'$CURRENTDATE'\.js/' ./templates/layout/base.gohtml
fi

docker compose run --rm nodejs bash -c "$GRUNTCMD"
