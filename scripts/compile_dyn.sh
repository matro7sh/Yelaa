#!/usr/bin/env bash

set -eou pipefail

#
# This script is a work-around to use gccgo with reflect2
# (which is a dependency of protobuf, so it's can of hard for us to
# fix this manually)
#
# Reference: https://github.com/modern-go/reflect2/issues/21
#
# I don't like this either, but this has been the only way to make proxychains
# work, because it relies on LD_PRELOAD to set proxies, so we should use gccgo
# to resolve libc dynamically
# https://github.com/golang/go/issues/31772#issuecomment-488322661
# https://github.com/Jguer/yay/issues/429#issuecomment-393661439
#

PKG_PATH="${GOPATH}/pkg/mod/github.com/modern-go/reflect2@v1.0.2"
FILENAME="unsafe_link.go"

PROGRAM_NAME=DynYelaa

compile () {
    echo "[+] Running compile with gccgo"
    go build -compiler gccgo -o ${PROGRAM_NAME}
}

backup_program_files () {
    echo "[+] Making backups of soon-to-be modified file"

    sudo cp -v "${PKG_PATH}/${FILENAME}" "/tmp/${FILENAME}"
    mv ~/.cache ~/.cache.bak
}

reset_cache () {

    echo "[+] Resetting file"
    sudo mv -v "/tmp/${FILENAME}" "${PKG_PATH}/${FILENAME}"

    mv ~/.cache.bak ~/.cache
}

replace_files () {
    echo "[+] Fixing reflect.unsafe_New call in ${PKG_PATH}/${FILENAME}"

    sudo sed -i 's/go:linkname unsafe_New reflect.unsafe_New/go:linkname unsafe_New reflect.unsafe__New/' "${PKG_PATH}/${FILENAME}"
    sudo sed -i 's/go:linkname unsafe_NewArray reflect.unsafe_NewArray/go:linkname unsafe_NewArray reflect.unsafe__NewArray/' "${PKG_PATH}/${FILENAME}"
}

backup_program_files

replace_files

reset_cache

echo "[+] Successfully generated ${PROGRAM_NAME}:"
file "${PROGRAM_NAME}"
