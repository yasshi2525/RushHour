#!/bin/bash

pushd "$(dirname "$0")/img" > /dev/null
outdir="$PWD/../../public/spritesheet"

find -mindepth 2 -type d  | while read line; do
    name="$(echo "$line" | cut -d "/" -f 2)"
    resolution="$(echo "$line" | cut -d "/" -f 3)"
    TexturePacker "$line" \
        --format pixijs4 \
        --data "$outdir/$name$resolution.json" \
        --extrude 0 \
        --algorithm Basic \
        --trim-mode None \
        --png-opt-level 0 \
        --max-width 4096 \
        --max-height 4096 \
        --disable-auto-alias
done