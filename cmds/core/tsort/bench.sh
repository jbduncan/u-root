#!/usr/bin/env bash

# Copyright 2012-2026 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# bench.sh runs Go and hyperfine benchmarks against two given Git-commit-based
# versions of tsort. The hyperfine benchmarks also compare these versions of
# tsort with the system tsort and uutils/coreutils' [2] tsort.
#
# This script accepts two arguments, each containing a Git commit pointing at
# two different versions of tsort.
#
# Usage
#
#     ./cmds/core/tsort/bench.sh <commit-before> <commit-after>
#
# Note: This script assumes that a system tsort and hyperfine [1] are installed
# and that uutils/coreutils is installed on the PATH with a "uu" prefix such
# their tsort is named `uutsort`.
#
# [1] https://github.com/sharkdp/hyperfine
# [2] https://github.com/uutils/coreutils

set -euo pipefail

current_branch_or_commit="$(git rev-parse --abbrev-ref HEAD)"
if [ "$current_branch_or_commit" = "HEAD" ]; then
    current_branch_or_commit="$(git rev-parse HEAD)"
fi
trap 'git checkout "$current_branch_or_commit"' EXIT

git checkout "$1"
go build -o ./tsort-before ./cmds/core/tsort
trap 'rm ./tsort-before; git checkout "$current_branch_or_commit"' EXIT
printf "\nRunning warmup Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=2 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-before...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/... | tee tsort-bench-before.txt

git checkout "$2"
go build -o ./tsort-after ./cmds/core/tsort
trap 'rm ./tsort-after; rm ./tsort-before; git checkout "$current_branch_or_commit"' EXIT
printf "\nRunning warmup Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=2 ./cmds/core/tsort/...
printf "\nRunning real Go benchmarks for ./tsort-after...\n"
go test -run=XXX -bench=Tsort -benchmem -count=10 ./cmds/core/tsort/... | tee tsort-bench-after.txt

go run golang.org/x/perf/cmd/benchstat@latest tsort-bench-before.txt tsort-bench-after.txt | tee tsort-bench-comparison.txt

graphsdir=$(go run ./cmds/core/tsort/gengraphs/)
trap 'rm -r "$graphsdir"; rm ./tsort-after; rm ./tsort-before; git checkout "$current_branch_or_commit"' EXIT

printf "\nRunning Hyperfine benchmarks for ./tsort-before, ./tsort-after, uutsort and tsort on acyclic graphs...\n"
for filepath in "$graphsdir"/*-acyclic-graph.txt; do
    markdown_file=$(basename "$filepath" .txt).md
    hyperfine --warmup 15 --runs 50 --shell=none --export-markdown="$markdown_file" \
        "./tsort-before $filepath" \
        "./tsort-after $filepath" \
        "uutsort $filepath" \
        "tsort $filepath"
done
printf "\nRunning Hyperfine benchmarks for ./tsort-before, ./tsort-after, uutsort and tsort on cyclic graphs...\n"
for filepath in "$graphsdir"/*-cyclic-graph.txt; do
    markdown_file=$(basename "$filepath" .txt).md
    hyperfine --warmup 15 --runs 50 --shell=none --export-markdown="$markdown_file" --ignore-failure \
        "./tsort-before $filepath" \
        "./tsort-after $filepath" \
        "uutsort $filepath" \
        "tsort $filepath"
done
