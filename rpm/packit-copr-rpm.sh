#!/usr/bin/env bash
PKG_NAME="prometheus-podman-exporter"
GIT_TOPDIR=$(git rev-parse --show-toplevel)
SPEC_FILE=${GIT_TOPDIR}/rpm/${PKG_NAME}.spec

VERSION=$(grep -E "^VERSION=" ${GIT_TOPDIR}/VERSION  | cut -d= -f2)
REVISION=$(grep -E "^REVISION=" ${GIT_TOPDIR}/VERSION  | cut -d= -f2)

echo $REVISION | grep 'dev'
if [ $? -eq 0 ] ; then
    VERSION="${VERSION}dev"
fi

sed -i "s/^Version:.*/Version: ${VERSION}/" $SPEC_FILE

git-archive-all -C "$GIT_TOPDIR" --prefix="${PKG_NAME}-${VERSION}/" "$GIT_TOPDIR/rpm/${PKG_NAME}-${VERSION}.tar.gz"

go mod vendor
tar czf "vendor-${VERSION}.tar.gz" vendor/
mv "vendor-${VERSION}.tar.gz" "$GIT_TOPDIR/rpm/"
