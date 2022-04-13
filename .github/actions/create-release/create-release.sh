#!/bin/bash -

export BOSH_VERSION
BOSH_VERSION=$(bosh --version | awk '{print $2}')
echo "BOSH Version: ${BOSH_VERSION}"
echo "::set-output name=bosh_version::${BOSH_VERSION}"

bosh_args=()
if [ -n "${INPUT_NAME}" ]; then
  bosh_args+=("--name" "${INPUT_NAME}")
fi
if [ -n "${INPUT_FORCE}" ]; then
  bosh_args+=("--force" "${INPUT_FORCE}")
fi
if [ "${INPUT_TIMESTAMP_VERSION}" = 'true' ]; then
  bosh_args+=("--timestamp-version")
fi
if [ -n "${INPUT_TARBALL}" ]; then
  bosh_args+=("--tarball" "${INPUT_TARBALL}")
fi
if [ "${INPUT_FINAL}" = 'true' ]; then
  bosh_args+=("--tarball")
fi
if [ -n "${INPUT_VERSION}" ]; then
  bosh_args+=("--version" "${INPUT_VERSION}")
fi
if [ -n "${INPUT_DIR}" ]; then
  bosh_args+=("--dir" "${INPUT_DIR}")
fi

bosh create-release "${bosh_args[@]}"
