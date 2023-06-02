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

CSSMIN=$(md5sum public/assets/css/reprogl.min.css | awk '{print $1}')
CSSMIN_NEW=$(md5sum public/assets/css/reprogl_temp.min.css | awk '{print $1}')

echo "reprogl.min.css      md5: $CSSMIN"
echo "reprogl_temp.min.css md5: $CSSMIN_NEW"

if [[ "$CSSMIN" != "$CSSMIN_NEW" ]]; then
  cat public/assets/css/reprogl_temp.min.css > public/assets/css/reprogl.min.css
  echo "Replace CSS file"
fi

if [[ "$GRUNTCMD" == "grunt" ]]; then
  JSMIN=$(md5sum public/assets/js/reprogl.min.js | awk '{print $1}')
  JSMIN_NEW=$(md5sum public/assets/js/reprogl_temp.min.js | awk '{print $1}')

  echo -e "\nreprogl.min.js       md5: $JSMIN"
  echo "reprogl_temp.min.js  md5: $JSMIN_NEW"

  if [[ "$JSMIN" != "$JSMIN_NEW" ]]; then
    cat public/assets/js/reprogl_temp.min.js > public/assets/js/reprogl.min.js
    echo "Replace JS file"
  fi
fi
