#!/bin/bash

script_dir="$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit 1
    pwd
)"
project_dir="$(
    cd "${script_dir}/.." || exit 1
    pwd
)"

if [[ -z ${APP} ]]; then
    APP="web"
fi

commit=${COMMIT:-"$(git rev-parse --short HEAD)"}
branch=${BRANCH:-"$(git rev-parse --abbrev-ref HEAD)"}
version=${branch}-${commit}

ldflags="-s -w"
ldflags+=" -X github.com/prometheus/common/version.Revision=${commit}"
ldflags+=" -X github.com/prometheus/common/version.Branch=${branch}"
ldflags+=" -X github.com/prometheus/common/version.Version=${version}"

GOOS=${OS:-'linux'} GOARCH='amd64' CGO_ENABLED=0 go build \
    -ldflags="${ldflags}" \
    -o "${project_dir}/bin/${APP}" \
    "${project_dir}/cmd/${APP}"
