#!/usr/bin/env bash
set -Eeuo pipefail

usage() {
    echo "Usage: bash tag-gplt-module.sh <module-name> <version>" >&2
    echo "Example: bash tag-gplt-module.sh math v1.0.0" >&2
    exit 2
}

[[ $# -eq 2 ]] || usage

module_name="$1"
version="$2"

# This wrapper accepts a single feature-module directory name.
if [[ ! "$module_name" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    echo "Error: provide a valid module name inside of pkg/gplt/modules" >&2
    exit 1
fi

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"

module_dir="pkg/gplt/modules/$module_name"

if [[ ! -f "$repo_root/$module_dir/go.mod" ]]; then
    echo "Error: module '$module_name' not found at '$module_dir'" >&2
    exit 1
fi

exec bash "$repo_root/release" "$module_dir" "$version"
