#!/usr/bin/env bash
#
# compile.sh
#
# Compiles each antas variant for the current platform ONLY -- no
# cross-compiling. antas-natif needs cgo + the native PDFium lib available
# via pkg-config on this machine (see the pdfium setup notes in the repo).
#
# Usage:
#   ./compile.sh
#
# Add more variants by adding a "name=tags" entry to BUILDS below, e.g.:
#   "antas-natif-turbo=natif pdfium_use_turbojpeg"

set -euo pipefail

# Each entry is "name=tags" (space-separated tags, or empty for none).
BUILDS=(
	"antas="
	"antas-turbo=pdfium_use_turbojpeg"
	"antas-natif=natif pdfium_use_turbojpeg" # Sounds misleading but antas-natif also use turbojpeg
)


# --- sanity checks -----------------------------------------------------

if ! command -v go >/dev/null 2>&1; then
	echo "Error: go is not installed or not on PATH." >&2
	exit 1
fi

needs_cgo=0
for entry in "${BUILDS[@]}"; do
	tags="${entry#*=}"
	if [[ " $tags " == *" natif "* ]]; then
		needs_cgo=1
	fi
done

if [ "$needs_cgo" -eq 1 ]; then
	if [ "$(go env CGO_ENABLED)" != "1" ]; then
		echo "Error: CGO_ENABLED is not 1 (currently '$(go env CGO_ENABLED)')," >&2
		echo "but a 'natif' build was requested. Run: export CGO_ENABLED=1" >&2
		exit 1
	fi
	if ! command -v pkg-config >/dev/null 2>&1; then
		echo "Error: pkg-config not found, needed to locate PDFium for 'natif' builds." >&2
		exit 1
	fi
	if ! pkg-config --exists pdfium 2>/dev/null; then
		echo "Error: pkg-config can't find 'pdfium'." >&2
		echo "Check PKG_CONFIG_PATH points at your pdfium.pc (see setup notes)." >&2
		exit 1
	fi
fi

echo "Building for $(go env GOOS)/$(go env GOARCH) (host platform, no cross-compiling)"

# --- build loop ----------------------------------------------------------

overall_status=0
built=()

for entry in "${BUILDS[@]}"; do
	name="${entry%%=*}"
	tags="${entry#*=}"

	echo "==> Building $name${tags:+ (tags: $tags)}"

	build_args=()
	if [ -n "$tags" ]; then
		build_args+=(-tags "$tags")
	fi
	build_args+=(-o "$name" .)

	if go build "${build_args[@]}"; then
		built+=("$name")
	else
		echo "Error: failed to build '$name'" >&2
		overall_status=1
	fi
done

if [ "$overall_status" -ne 0 ]; then
	echo "One or more builds failed. See errors above." >&2
	exit "$overall_status"
fi

echo "Done. Built: ${built[*]}"