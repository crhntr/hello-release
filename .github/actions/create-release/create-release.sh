#!/bin/bash -l

export BOSH_VERSION
BOSH_VERSION=$(bosh --version | awk '{print $2}')

env

echo "BOSH Version: ${BOSH_VERSION}"

echo "::set-output name=bosh_version::${BOSH_VERSION}"