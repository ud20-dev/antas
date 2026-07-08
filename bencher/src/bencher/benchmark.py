#!/usr/bin/env python3
"""
plot_benchmark.py

Turn hyperfine JSON benchmark results (produced by benchmark.sh) into a
grouped bar chart comparing py-antas vs go-antas across every PDF sample.

Install:
    pip install matplotlib numpy

Usage:
    python plot_benchmark.py [results_dir] [output.png] [samples_dir]

Defaults:
    results_dir = benchmark_results
    output.png  = benchmark_comparison.png
    samples_dir = tests_samples
"""

from __future__ import annotations

import json
import sys
from datetime import datetime
from pathlib import Path

try:
    import numpy as np
except ImportError as e:
    print(f"Error: missing dependency 'numpy'. Install it with: pip install numpy\n({e})")
    sys.exit(1)

try:
    import matplotlib.pyplot as plt
except ImportError as e:
    print(f"Error: missing dependency 'matplotlib'. Install it with: pip install matplotlib\n({e})")
    sys.exit(1)

try:
    import pypdfium2 as pdfium
except ImportError as e:
    print(f"Error: missing dependency 'pypdfium2'. Install it with: pip install pypdfium2\n({e})")
    sys.exit(1)


def nice_list(items: list[str]) -> str:
    if not items:
        return ""
    if len(items) == 1:
        return items[0]
    elif len(items) == 2:
        return items[0] + " vs " + items[1]
    return ", ".join(items[:-1]) + " and " + items[-1]

def get_page_count(pdf_path: Path) -> int | None:
    """Returns the page count of a PDF, or None if it can't be read."""
    if not pdf_path.is_file():
        return None
    try:
        pdf = pdfium.PdfDocument(pdf_path)
        try:
            return len(pdf)
        finally:
            pdf.close()
    except Exception:
        return None


def load_results(results_dir: Path) -> dict[str, dict[str, tuple[float, float]]]:
    """
    Reads every hyperfine JSON export in `results_dir` and returns:
        {sample_filename: {command_name: (mean_seconds, stddev_seconds)}}
    """
    data: dict[str, dict[str, tuple[float, float]]] = {}

    json_files = sorted(results_dir.glob("*.json"))
    if not json_files:
        raise FileNotFoundError(f"no .json files found in '{results_dir}'")

    for json_path in json_files:
        try:
            with json_path.open() as f:
                payload = json.load(f)
        except Exception as e:
            raise ValueError(f"failed to parse '{json_path}': {e}") from e

        sample_name = json_path.stem
        data[sample_name] = {}

        for result in payload.get("results", []):
            command_name = result.get("command", "unknown")
            mean = result.get("mean", 0.0)
            stddev = result.get("stddev", 0.0)
            data[sample_name][command_name] = (mean, stddev)

    return data


def build_page_count_labels(samples: list[str], samples_dir: Path) -> list[str]:
    """
    Builds x-axis labels combining each sample's name and page count,
    e.g. "cv (3 Pages)". Falls back to just the filename if the PDF
    can't be found or read, so a missing/unreadable file doesn't break
    the whole plot.
    """
    labels = []
    for sample in samples:
        pdf_path = samples_dir / f"{sample}.pdf"
        page_count = get_page_count(pdf_path)
        if page_count is not None:
            unit = "Page" if page_count == 1 else "Pages"
            labels.append(f"{sample} ({page_count} {unit})")
        else:
            labels.append(sample)  # fallback: couldn't read page count
    return labels


def plot_comparison(
    data: dict[str, dict[str, tuple[float, float]]],
    output_path: Path,
    samples_dir: Path,
) -> None:
    samples = list(data.keys())
    commands = sorted({cmd for per_file in data.values() for cmd in per_file})

    if not commands:
        raise ValueError("no commands found in results")

    labels = build_page_count_labels(samples, samples_dir)

    x = np.arange(len(samples))
    width = 0.8 / len(commands)

    fig, ax = plt.subplots(figsize=(max(8, len(samples) * 1.2), 6))

    for i, command in enumerate(commands):
        means = []
        errs = []
        for sample in samples:
            mean, stddev = data[sample].get(command, (0.0, 0.0))
            means.append(mean)
            errs.append(stddev)

        offset = (i - (len(commands) - 1) / 2) * width
        ax.bar(x + offset, means, width, yerr=errs, capsize=4, label=command)

    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M")
    title = (
        f"{nice_list(commands)} — {len(samples)} PDF sample"
        f"{'s' if len(samples) != 1 else ''} — {timestamp}"
    )

    ax.set_xlabel("PDF sample")
    ax.set_ylabel("Time (seconds)")
    ax.set_title(title)
    ax.set_xticks(x)
    ax.set_xticklabels(labels, rotation=45, ha="right")
    ax.legend()
    ax.grid(axis="y", linestyle="--", alpha=0.4)

    fig.tight_layout()
    fig.savefig(output_path, dpi=150)
    print(f"Saved chart to '{output_path}'")


def main() -> None:
    results_dir = Path(sys.argv[1]) if len(sys.argv) > 1 else Path("benchmark_results")
    output_path = Path(sys.argv[2]) if len(sys.argv) > 2 else Path("benchmark_comparison.png")
    samples_dir = Path(sys.argv[3]) if len(sys.argv) > 3 else Path("tests_samples")

    if not results_dir.is_dir():
        print(f"Error: results directory '{results_dir}' does not exist.")
        sys.exit(1)

    try:
        data = load_results(results_dir)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

    try:
        plot_comparison(data, output_path, samples_dir)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()