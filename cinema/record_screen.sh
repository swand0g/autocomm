#!/bin/bash

DIRECTORY="$(dirname "$0")"

function setup_binary() {
  TEMP_DIR="/tmp/autocomm-$(date '+%s')"
  mkdir "$TEMP_DIR"
  cp ./target/release/autocomm "$TEMP_DIR/autocomm"
  export PATH=$TEMP_DIR:$PATH
  export autocomm_DIR=$TEMP_DIR/.autocomm

  # First run of the binary might be slower due to anti-virus software
  echo "Using $(which autocomm)"
  echo "  with version $(autocomm --version)"
}

setup_binary

RECORDING_PATH=$DIRECTORY/screen_recording

(rm -rf "$RECORDING_PATH" &> /dev/null || true)

asciinema rec -c "$DIRECTORY/recorded_screen_script.sh" "$RECORDING_PATH"
sed "s@$TEMP_DIR@~@g" "$RECORDING_PATH" | \
  svg-term                      \
    --window                    \
    --out "cinema/autocomm.svg" \
    --height=17                 \
    --width=70