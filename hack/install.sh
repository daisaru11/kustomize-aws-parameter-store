#!/usr/bin/env bash

set -e

help() {
    cat <<'EOF'
Install a binary release of daisaru11/kustomize-aws-parameter-store
Usage:
    install.sh [options]
Options:
    -h, --help      Display this message
    -f, --force     Force overwriting an existing binary
    --tag TAG       Tag (version) of the binary to install (default <latest release>)
    --target TARGET Install the release compiled for $TARGET
    --to LOCATION   Where to install the binary (default $XDG_CONFIG_HOME/kustomize/plugin/kustomize.daisaru11.dev/v1/awsparameterstoresecret)
EOF
}

_log() {
    echo "install.sh: $1"
}

_log_err() {
    _log "$1" >&2
}

require() {
    if ! command -v $1 > /dev/null 2>&1; then
        err "require $1 (command not found)"
    fi
}

force=false
while test $# -gt 0; do
    case $1 in
        --force | -f)
            force=true
            ;;
        --help | -h)
            help
            exit 0
            ;;
        --tag)
            tag=$2
            shift
            ;;
        --target)
            target=$2
            shift
            ;;
        --to)
            dest=$2
            shift
            ;;
        *)
            ;;
    esac
    shift
done

# Dependencies
require basename
require curl
require install
require mkdir
require mktemp
require tar

# Optional dependencies
if [ -z $tag ] || [ -z $target ]; then
    require cut
fi

if [ -z $tag ]; then
    require rev
fi

if [ -z $target ]; then
    require grep
fi

git="daisaru11/kustomize-aws-parameter-store"

url="https://github.com/$git"
_log_err "GitHub repository: $url"

url="$url/releases"

if [ -z $tag ]; then
  tag=$(curl -s "$url/latest" | cut -d'"' -f2 | rev | cut -d'/' -f1 | rev)
  _log_err "Tag: latest ($tag)"
else
  _log_err "Tag: $tag"
fi

if [ -z $target ]; then
  target=linux
  if [[ "$OSTYPE" == "darwin"* ]]; then
    target=darwin
  fi
fi

_log_err "Target: $target"

if [ -z $dest ]; then
  config_home=${XDG_CONFIG_HOME:-$HOME/.config}
  dest="$config_home/kustomize/plugin/kustomize.daisaru11.dev/v1/awsparameterstoresecret/AWSParameterStoreSecret"
fi

if [ -e "$dest" ] && [ $force = false ]; then
  _log_err "the binary already exists in $dest"
  exit 1
fi

_log_err "Installing to: $dest"

url="$url/download/$tag/AWSParameterStoreSecret-$target-amd64"

mkdir -p $(dirname $dest)

curl -sL -o $dest $url 
chmod a+x $dest
