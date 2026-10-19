#!/usr/bin/env bash
set -Eeuo pipefail

usage() {
    echo "Usage: bash tag-go-module.sh <module-directory> <version>" >&2
    echo "Example: bash tag-go-module.sh path/to/module v1.0.0" >&2
    exit 2
}

[[ $# -eq 2 ]] || usage

module_dir="${1#./}"
version="$2"

# Accept semantic versions such as v1.0.0 and v1.2.0-rc.1.
if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
    echo "Error: invalid version '$version'; expected a version like v1.0.0" >&2
    exit 1
fi

# Resolve the repository from this script's location.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"

# Require a repository-relative module directory.
case "$module_dir" in
    ""|/*|..|../*|*/../*|*/..)
        echo "Error: provide a valid repository-relative module directory" >&2
        exit 1
        ;;
esac

module_path="$repo_root/$module_dir"

if [[ ! -f "$module_path/go.mod" ]]; then
    echo "Error: no go.mod found in '$module_dir'" >&2
    exit 1
fi

# Confirm this is a separate module within the same repository.
module_root="$(git -C "$module_path" rev-parse --show-toplevel)"
if [[ "$module_root" != "$repo_root" ]]; then
    echo "Error: module is not inside the expected Git repository" >&2
    exit 1
fi

tag="$module_dir/$version"

# Require committed changes before releasing.
if [[ -n "$(git -C "$repo_root" status --porcelain)" ]]; then
    echo "Error: working tree is not clean. Commit or stash your changes first." >&2
    exit 1
fi

# Never overwrite an existing local tag.
if git -C "$repo_root" rev-parse -q --verify "refs/tags/$tag" >/dev/null; then
    echo "Error: tag '$tag' already exists locally" >&2
    exit 1
fi

# Check the remote before creating the tag.
if ! git -C "$repo_root" remote get-url origin >/dev/null 2>&1; then
    echo "Error: Git remote 'origin' is not configured" >&2
    exit 1
fi

if git -C "$repo_root" ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1; then
    echo "Error: tag '$tag' already exists on origin" >&2
    exit 1
fi

echo "Testing module: $module_dir"
(
    cd "$module_path"
    go test ./...
)

echo "Creating annotated tag: $tag"
git -C "$repo_root" tag -a "$tag" -m "tag-go-module $tag"

echo "Pushing tag to origin..."
if ! git -C "$repo_root" push origin "refs/tags/$tag"; then
    echo "Error: push failed. The local tag '$tag' remains; resolve the issue before retrying." >&2
    exit 1
fi

echo "tag-go-moduled $tag"
