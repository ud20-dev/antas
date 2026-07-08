#!/usr/bin/env bash
#
# benchmark.sh
#
# For every PDF in tests_samples/, benchmark all configured implementations
# against each other with hyperfine (3 warmup runs each).
#
# Usage:
#   ./benchmark.sh
#
# Requirements:
#   - hyperfine installed (https://github.com/sharkdp/hyperfine)
#   - each binary listed in BINARIES built and present next to this script
#     (or edit BINARIES below to point elsewhere)

set -euo pipefail

SAMPLES_DIR="tests_samples"
RESULTS_DIR="benchmark_results"

# Each entry is "name=path". Add/remove entries here to change what gets
# benchmarked -- no other part of the script needs to change.
BINARIES=(
	"antas-turbo=./antas-turbo"
	"antas-natif=./antas-natif"
	"antas=./antas"
)

# --- sanity checks -----------------------------------------------------

if ! command -v hyperfine >/dev/null 2>&1; then
	echo "Error: hyperfine is not installed or not on PATH." >&2
	echo "Install instructions: https://github.com/sharkdp/hyperfine#installation" >&2
	exit 1
fi

if [ ${#BINARIES[@]} -eq 0 ]; then
	echo "Error: BINARIES array is empty. Add at least one \"name=path\" entry." >&2
	exit 1
fi

names=()
paths=()

for entry in "${BINARIES[@]}"; do
	name="${entry%%=*}"
	path="${entry#*=}"

	if [ -z "$name" ] || [ -z "$path" ]; then
		echo "Error: malformed BINARIES entry '$entry' (expected \"name=path\")." >&2
		exit 1
	fi

	if [ ! -x "$path" ]; then
		echo "Error: '$path' (for '$name') not found or not executable." >&2
		echo "Build it first, e.g.: go build -o ${path#./} ." >&2
		exit 1
	fi

	names+=("$name")
	paths+=("$path")
done

if [ ! -d "$SAMPLES_DIR" ]; then
	echo "Error: samples directory '$SAMPLES_DIR' does not exist." >&2
	exit 1
fi

shopt -s nullglob
pdf_files=("$SAMPLES_DIR"/*.pdf)
shopt -u nullglob

if [ ${#pdf_files[@]} -eq 0 ]; then
	echo "Error: no .pdf files found in '$SAMPLES_DIR'." >&2
	exit 1
fi

mkdir -p "$RESULTS_DIR"

# --- benchmark loop ------------------------------------------------------

overall_status=0

for file in "${pdf_files[@]}"; do
	filename="$(basename "$file")"
	result_json="$RESULTS_DIR/${filename%.pdf}.json"
	result_md="$RESULTS_DIR/${filename%.pdf}.md"

	echo "==> Benchmarking $filename"

	hyperfine_args=(
		--warmup 3
		--export-json "$result_json"
		--export-markdown "$result_md"
	)

	for i in "${!names[@]}"; do
		hyperfine_args+=(--command-name "${names[$i]}" "${paths[$i]} \"$file\"")
	done

	if ! hyperfine "${hyperfine_args[@]}"; then
		echo "Error: hyperfine failed on '$filename'" >&2
		overall_status=1
		continue
	fi
done

if [ "$overall_status" -ne 0 ]; then
	echo "One or more benchmarks failed. See errors above." >&2
	exit "$overall_status"
fi

echo "Done. Results saved in '$RESULTS_DIR/'."