#!/usr/bin/env python3

import argparse
import re
import subprocess
from pathlib import Path


def run_and_tee(cmd, output_path):
    with output_path.open("w") as output:
        process = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
        )
        assert process.stdout is not None
        for line in process.stdout:
            print(line, end="")
            output.write(line)
        returncode = process.wait()
        if returncode != 0:
            raise subprocess.CalledProcessError(returncode, cmd)


def main():
    parser = argparse.ArgumentParser(
        description="Run and compare benchmarks for a Go package."
    )
    parser.add_argument("folder", help="package folder to benchmark, for example lex")
    args = parser.parse_args()

    folder = Path(args.folder)
    if not folder.is_dir():
        parser.error(f"folder does not exist: {args.folder}")

    benchmark_dir = Path("benchmarks") / folder.name
    benchmark_dir.mkdir(parents=True, exist_ok=True)

    existing = []
    for path in benchmark_dir.glob("*.bench"):
        match = re.fullmatch(r"(\d+)-.+\.bench", path.name)
        if match:
            existing.append((int(match.group(1)), path))

    previous = max(existing, default=None, key=lambda item: item[0])
    number = previous[0] + 1 if previous else 1

    commit = subprocess.check_output(
        ["git", "rev-parse", "--short", "HEAD"], text=True
    ).strip()
    benchmark = benchmark_dir / f"{number:02d}-{commit}.bench"

    package = f"./{folder.as_posix()}"
    run_and_tee(
        ["go", "test", "-bench=.", "-benchmem", "-benchtime", "1s", package],
        benchmark,
    )

    print(f"Saved benchmark to {benchmark}")

    if previous:
        comparison = benchmark_dir / f"{previous[0]:02d}vs{number:02d}.benchstat"
        run_and_tee(["benchstat", str(previous[1]), str(benchmark)], comparison)
        print(f"Saved comparison to {comparison}")


if __name__ == "__main__":
    main()
