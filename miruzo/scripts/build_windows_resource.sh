#!/usr/bin/env bash
set -eu

os="${1:?os is required}"

if [ "$os" != "windows" ]; then
	exit 0
fi

arch="${2:?arch is required}"
version="${3:?version is required}"
raw_version="${4:?raw_version is required}"

machine_flag=""
case "$arch" in
	amd64)
		machine_flag="-64"
		;;
	386)
		machine_flag=""
		;;
	arm)
		machine_flag="-arm"
		;;
	arm64)
		machine_flag="-arm64"
		;;
	*)
		echo "unsupported arch: $arch" >&2
		exit 1
		;;
esac

file_version_flag=" \
	-ver-major=$(echo "$raw_version" | cut -d. -f1) \
	-ver-minor=$(echo "$raw_version" | cut -d. -f2) \
	-ver-patch=$(echo "$raw_version" | cut -d. -f3) \
	-ver-build=0 \
	"

product_version_flag=" \
	-product-ver-major=$(echo "$raw_version" | cut -d. -f1) \
	-product-ver-minor=$(echo "$raw_version" | cut -d. -f2) \
	-product-ver-patch=$(echo "$raw_version" | cut -d. -f3) \
	-product-ver-build=0 \
	"

# miruzo-api
goversioninfo \
	-description="miruzo API" \
	-file-version="$raw_version.0" \
	-internal-name="miruzo-api" \
	-manifest=./assets/app.manifest \
	-o=cmd/miruzo-api/resource_windows_${arch}.syso \
	-original-name="miruzo-api.exe" \
	-product-version="$version" \
	$machine_flag \
	$file_version_flag \
	$product_version_flag \
	assets/versioninfo.json

# miruzo-cli
goversioninfo \
	-description="miruzo CLI" \
	-file-version="$raw_version.0" \
	-internal-name="miruzo-cli" \
	-manifest=assets/app.manifest \
	-o=cmd/miruzo-cli/resource_windows_${arch}.syso \
	-original-name="miruzo-cli.exe" \
	-product-version="$version" \
	$machine_flag \
	$file_version_flag \
	$product_version_flag \
	assets/versioninfo.json
