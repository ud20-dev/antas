#!/usr/bin/env pwsh
#
# benchmark.ps1
#
# For every PDF in tests_samples/, benchmark all configured implementations
# against each other with hyperfine (3 warmup runs each).
#
# Usage:
#   ./benchmark.ps1
#
# Requirements:
#   - hyperfine installed (https://github.com/sharkdp/hyperfine)
#   - each binary listed in $Binaries built and present next to this script
#     (or edit $Binaries below to point elsewhere)

$ErrorActionPreference = "Stop"

$SamplesDir = "tests_samples"
$ResultsDir = "benchmark_results"

# Each entry is a hashtable @{ Name = ...; Path = ... }. Add/remove entries
# here to change what gets benchmarked -- no other part of the script needs
# to change.
$Binaries = @(
    @{ Name = "antas-turbo"; Path = ".\antas-turbo.exe" }
    @{ Name = "antas-natif"; Path = ".\antas-natif.exe" }
    @{ Name = "antas";       Path = ".\antas.exe" }
)

# --- sanity checks -----------------------------------------------------

if (-not (Get-Command hyperfine -ErrorAction SilentlyContinue)) {
    Write-Error "hyperfine is not installed or not on PATH.`nInstall instructions: https://github.com/sharkdp/hyperfine#installation"
    exit 1
}

if ($Binaries.Count -eq 0) {
    Write-Error "Binaries array is empty. Add at least one Name/Path entry."
    exit 1
}

$names = @()
$paths = @()

foreach ($entry in $Binaries) {
    $name = $entry.Name
    $path = $entry.Path

    if ([string]::IsNullOrWhiteSpace($name) -or [string]::IsNullOrWhiteSpace($path)) {
        Write-Error "Malformed Binaries entry (expected Name/Path): $($entry | Out-String)"
        exit 1
    }

    if (-not (Test-Path -Path $path -PathType Leaf)) {
        Write-Error "'$path' (for '$name') not found.`nBuild it first, e.g.: go build -o $($path.TrimStart('.','\')) ."
        exit 1
    }

    $names += $name
    $paths += $path
}

if (-not (Test-Path -Path $SamplesDir -PathType Container)) {
    Write-Error "Samples directory '$SamplesDir' does not exist."
    exit 1
}

$pdfFiles = Get-ChildItem -Path $SamplesDir -Filter "*.pdf" -File

if ($pdfFiles.Count -eq 0) {
    Write-Error "No .pdf files found in '$SamplesDir'."
    exit 1
}

New-Item -ItemType Directory -Force -Path $ResultsDir | Out-Null

# --- benchmark loop ------------------------------------------------------

$overallStatus = 0

foreach ($file in $pdfFiles) {
    $filename = $file.Name
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($filename)
    $resultJson = Join-Path $ResultsDir "$baseName.json"
    $resultMd = Join-Path $ResultsDir "$baseName.md"

    Write-Host "==> Benchmarking $filename"

    $hyperfineArgs = @(
        "--warmup", "3",
        "--export-json", $resultJson,
        "--export-markdown", $resultMd
    )

    for ($i = 0; $i -lt $names.Count; $i++) {
        $hyperfineArgs += "--command-name"
        $hyperfineArgs += $names[$i]
        $hyperfineArgs += "$($paths[$i]) `"$($file.FullName)`""
    }

    & hyperfine @hyperfineArgs
    if ($LASTEXITCODE -ne 0) {
        Write-Error "hyperfine failed on '$filename'"
        $overallStatus = 1
        continue
    }
}

if ($overallStatus -ne 0) {
    Write-Error "One or more benchmarks failed. See errors above."
    exit $overallStatus
}

Write-Host "Done. Results saved in '$ResultsDir/'."