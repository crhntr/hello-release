#!/bin/bash

set -euo pipefail

cd "${GITHUB_WORKSPACE}" || exit 1

git config --global --add safe.directory "${GITHUB_WORKSPACE}"
git checkout -b create-release

ls -la

export BOSH_VERSION
BOSH_VERSION=$(bosh --version | awk '{print $2}')
echo "BOSH Version: ${BOSH_VERSION}"
echo "bosh_version=${BOSH_VERSION}" >> "${GITHUB_OUTPUT}"

bosh_args=()
if [ -n "${INPUT_NAME}" ]; then
  bosh_args+=("--name" "${INPUT_NAME}")
fi
if [ -n "${INPUT_FORCE}" ]; then
  bosh_args+=("--force" "${INPUT_FORCE}")
fi
if [ "${INPUT_TIMESTAMP_VERSION}" = 'true' ]; then
  printf "INPUT_TIMESTAMP_VERSION not supported"
  exit 1
fi
if [ -n "${INPUT_TARBALL}" ]; then
  bosh_args+=("--tarball" "${INPUT_TARBALL}")
fi
if [ "${INPUT_FINAL}" = 'true' ]; then
  bosh_args+=("--final")
fi

export BOSH_RELEASE_VERSION
if [ -n "${INPUT_VERSION}" ]; then
  BOSH_RELEASE_VERSION="${INPUT_VERSION#"v"}"
  bosh_args+=("--version" "${BOSH_RELEASE_VERSION}")
elif [ "${INPUT_TIMESTAMP_VERSION}" = 'true' ]; then
   printf "required variable INPUT_VERSION not set"
  exit 1
fi

if [ -n "${INPUT_DIR}" ]; then
  bosh_args+=("--dir" "${INPUT_DIR}")
fi

export BOSH_SHA2=""
if [ "${INPUT_SHA2}" = 'false' ]; then
  unset BOSH_SHA2
fi

export GOOGLE_APPLICATION_CREDENTIALS=/tmp/key.yml
echo "${CONFIG_PRIVATE_JSON_KEY}" > "${GOOGLE_APPLICATION_CREDENTIALS}"

bosh --non-interactive create-release "${bosh_args[@]}"

rm "${GOOGLE_APPLICATION_CREDENTIALS}"

if [ -f "${INPUT_TARBALL}" ]; then
  chmod 666 "${INPUT_TARBALL}"
fi

if [ "${INPUT_FINAL}" = 'true' ]; then
  git config --global user.name "${GIT_USER_NAME}"
  git config --global user.email "${GIT_USER_EMAIL}"
  git add .
  set -x
  git commit -am "Create Release ${BOSH_RELEASE_VERSION}"
  git fetch origin main
  git switch main
  git pull --ff-only origin main
  git switch create-release
  git rebase main
  git switch main
  git merge --ff-only create-release
  git tag "v${INPUT_VERSION}"
  git push origin main "v${INPUT_VERSION}"
fi

if [[ "${INPUT_VERSION}" =~ "-" ]]; then
  echo "pre-release-version=true" >> "${GITHUB_OUTPUT}"
else
  echo "pre-release-version=false" >> "${GITHUB_OUTPUT}"
fi