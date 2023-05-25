#!/usr/bin/env bash

GRUNTCMD="grunt"

while getopts ":s" opt; do
  case $opt in
    s)
      GRUNTCMD="$GRUNTCMD style"
      ;;
  esac
done

docker compose run --rm nodejs bash -c "$GRUNTCMD"
