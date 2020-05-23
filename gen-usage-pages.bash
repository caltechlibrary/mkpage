#!/bin/bash

if [ ! -d bin ]; then
    echo "Run "make" before running this script."
    exit
fi
ls -1 bin/ | while read ITEM; do
   D=$(basename "${ITEM}")
   mkdir -p "docs/${D}"
   "bin/${ITEM}" -generate-markdown > "docs/${D}/index.md"
done
