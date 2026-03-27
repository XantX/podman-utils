#!/bin/bash

set -e

REPO="podutil/podutil"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="latest"

if [ "$1" != "" ]; then
    VERSION="$1"
fi

detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "Arquitectura no soportada: $ARCH"; exit 1 ;;
    esac
}

download() {
    detect_os

    local format="tar.gz"
    if [ "$OS" = "windows" ]; then
        format="zip"
    fi

    if [ "$VERSION" = "latest" ]; then
        local download_url="https://github.com/${REPO}/releases/${VERSION}/download/podutil_${VERSION}_${OS}_${ARCH}.${format}"
    else
        local download_url="https://github.com/${REPO}/releases/download/v${VERSION}/podutil_v${VERSION}_${OS}_${ARCH}.${format}"
    fi

    echo "Descargando podutil v${VERSION} para ${OS}/${ARCH}..."

    local tmp_dir=$(mktemp -d)
    cd "$tmp_dir"

    if command -v curl &> /dev/null; then
        curl -sL "$download_url" -o podutil.${format}
    elif command -v wget &> /dev/null; then
        wget -q "$download_url" -O podutil.${format}
    else
        echo "Error: Se requiere curl o wget"
        exit 1
    fi

    # Extraer
    if [ "$format" = "tar.gz" ]; then
        tar -xzf podutil.${format}
    else
        unzip -q podutil.${format}
    fi

    # Instalar
    mkdir -p "$INSTALL_DIR"
    if [ -f "podutil" ]; then
        mv podutil "$INSTALL_DIR/podutil"
        chmod +x "$INSTALL_DIR/podutil"
    elif [ -f "podutil.exe" ]; then
        mv podutil.exe "$INSTALL_DIR/podutil.exe"
        chmod +x "$INSTALL_DIR/podutil.exe"
    fi

    echo "Instalado en: $INSTALL_DIR/podutil"
    echo "Asegurate de agregar $INSTALL_DIR al PATH:"
    echo "  export PATH=\$PATH:$INSTALL_DIR"

    cd -
    rm -rf "$tmp_dir"
}

download

echo "Instalación completada!"