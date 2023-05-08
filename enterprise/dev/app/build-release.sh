#!/usr/bin/env bash

set -eu

GCLOUD_APP_CREDENTIALS_FILE=${GCLOUD_APP_CREDENTIALS_FILE-$HOME/.config/gcloud/application_default_credentials.json}
cd "$(dirname "${BASH_SOURCE[0]}")"/../../.. || exit 1

bazel_build() {
  local bazel_cmd
  local platform
  platform=$1
  bazel_cmd="bazel"
  if [[ ${GITHUB_ACTIONS:-""} == "true" ]]; then
    bazel_cmd="${bazel_cmd} --bazelrc=.aspect/bazelrc/github.bazelrc"
  fi

  echo "[Bazel] Building Sourcegraph Backend (${VERSION}) for platform: ${platform}"
  ${bazel_cmd} build //enterprise/cmd/sourcegraph:sourcegraph --stamp --workspace_status_command=./enterprise/dev/app/app_stamp_vars.sh

  out=$(bazel cquery //enterprise/cmd/sourcegraph:sourcegraph --output=files)
  mkdir -p ".bin"
  cp -vf "${out}" .bin/sourcegraph-backend-${platform}
}

create_version() {
    local sha=$(git rev-parse --short HEAD)
    if [[ -z ${sha} ]]; then
      sha=${GITHUB_SHA:-""}
    fi

    local build="insiders"
    if [[ ${RELEASE_BUILD} == 1 ]]; then
      build="release"
    fi
    echo "$(date '+%Y.%-m.%-d')+${build}.${sha}"
}

set_version() {
  if [[ ${CI:-""} == "true" ]]; then
    # Backend version
    export VERSION=$(cat ./src-tauri/VERSION)
    export APP_VERSION=$(create_version)
  else
    # Backend version
    export VERSION=$(git rev-parse HEAD)
    export APP_VERSION="0.0.0+dev.fakesha"
  fi

  tmp=$(mktemp)
  jq --arg version ${APP_VERSION} '.package.version = $version' ./src-tauri/tauri.conf.json > ${tmp}
  mv ${tmp} ./src-tauri/tauri.conf.json
}

set_platform() {
  # We need to determine the platform string for the sourcegraph-backend binary
  local arch=""
  local platform=""
  case "$(uname -m)" in
    "amd64")
      arch="x86_64"
      ;;
    "arm64")
      arch="aarch64"
      ;;
    *)
      arch=$(uname -m)
  esac

  case "$(uname -s)" in
    "Darwin")
      platform="${arch}-apple-darwin"
      ;;
    "Linux")
      platform="${arch}-unknown-linux-gnu"
      ;;
    *)
      platform="${arch}-unknown-unknown"
  esac

  export PLATFORM=${platform}
}

set_platform
set_version
bazel_build ${PLATFORM}
echo "[Tauri] Building Application (${APP_VERSION})"]
NODE_ENV=production pnpm run build-app-shell
pnpm tauri build
