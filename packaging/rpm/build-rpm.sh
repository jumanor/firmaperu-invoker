#!/bin/bash

# Script para construir el paquete RPM de Firma Peru Invoker

# Verificamos si rpmbuild está instalado
if ! command -v rpmbuild &> /dev/null; then
    echo "Error: rpmbuild no está instalado."
    echo "Para instalar rpmbuild, ejecute: sudo yum install rpm-build rpmdevtools"
    exit 1
fi

# Definimos directorios
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
BUILD_DIR="/tmp/firmaperu-invoker-rpm-build"
VERSION=$(cat "$ROOT_DIR/version.txt" 2>/dev/null || echo "1.5.1")

# Creamos estructura de directorios para rpmbuild
mkdir -p $BUILD_DIR/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Copiamos archivos necesarios
echo "Copiando archivos al directorio sources..."

cp "$ROOT_DIR/main" "$BUILD_DIR/SOURCES/"
cp "$ROOT_DIR/config.properties.example" "$BUILD_DIR/SOURCES/"
cp -r "$ROOT_DIR/public" "$BUILD_DIR/SOURCES/"

cp "$SCRIPT_DIR/firmaperu-invoker.service" "$BUILD_DIR/SOURCES/"
cp "$SCRIPT_DIR/firmaperu-invoker.spec" "$BUILD_DIR/SPECS/"

# Actualizamos la versión en el spec file
sed -i "s/%define version .*/%define version $VERSION/" "$BUILD_DIR/SPECS/firmaperu-invoker.spec"

# Construimos el RPM
echo "Construyendo el paquete RPM..."
rpmbuild --define "_topdir $BUILD_DIR" -bb "$BUILD_DIR/SPECS/firmaperu-invoker.spec"

# Movemos el RPM generado al directorio actual
mkdir -p "$SCRIPT_DIR/dist"
find "$BUILD_DIR/RPMS" -name "*.rpm" -exec cp {} "$SCRIPT_DIR/dist/" \;

echo "El paquete RPM se ha creado en: $SCRIPT_DIR/dist/"
