#!/usr/bin/env bash
set -euo pipefail

APP_NAME=orchastration
VERSION=${VERSION:-dev}

mkdir -p dist/package/linux

cp dist/${APP_NAME} dist/package/linux/
cp configs/config.example.toml dist/package/linux/${APP_NAME}.toml

( cd dist/package && tar -czf ../${APP_NAME}-${VERSION}-linux-amd64.tar.gz linux )

if command -v dpkg-deb >/dev/null 2>&1; then
  PKGDIR=dist/package/deb
  mkdir -p ${PKGDIR}/DEBIAN
  mkdir -p ${PKGDIR}/usr/local/bin
  mkdir -p ${PKGDIR}/etc/${APP_NAME}

  cat > ${PKGDIR}/DEBIAN/control <<CONTROL
Package: ${APP_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: ${APP_NAME} team
Description: Cross-platform orchestration CLI
CONTROL

  cp dist/${APP_NAME} ${PKGDIR}/usr/local/bin/${APP_NAME}
  cp configs/config.example.toml ${PKGDIR}/etc/${APP_NAME}/config.toml

  dpkg-deb --build ${PKGDIR} dist/${APP_NAME}_${VERSION}_amd64.deb
fi
