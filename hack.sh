#!/usr/bin/env bash

# Configuration

: ${BOSH_RELEASE_VERSION?"Need to set BOSH_RELEASE_VERSION"}
: ${BOSH_DEPLOYMENT?"Need to set BOSH_DEPLOYMENT"}
: ${STEMCELL_VERSION?"Need to set STEMCELL_VERSION"}
: ${TILE_VERSION?"Need to set TILE_VERSION"}
export TILE_STEMCELL_VERSION="${TILE_STEMCELL_VERSION:="${STEMCELL_VERSION}"}"

bosh upload-release --version="${BOSH_RELEASE_VERSION}"

bosh upload-stemcell \
  "https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-xenial-go_agent?v=${STEMCELL_VERSION}"

yq -y ".stemcells[0].version = ${STEMCELL_VERSION}" \
  ./examples/deployment/manifest.yml | sponge ./examples/deployment/manifest.yml

bosh deploy ./examples/deployment/manifest.yml

rm -f examples/tile/release/hello-release-*.pivotal

bosh export-release \
  "hello-release/0.11.0" \
  "ubuntu-xenial/621.23" \
  --dir=examples/tile/releases

# Deploy tile with release compiled with INITIAL_STEMCELL

cd examples/tile || exit 1
  yq -y ".stemcell_criteria.version = ${TILE_STEMCELL_VERSION}" \
    Kilnfile.lock | sponge Kilnfile.lock

  kiln bake --version="${TILE_VERSION}"

  curl "https://bosh-gce-light-stemcells.s3-accelerate.amazonaws.com/${STEMCELL_VERSION}/light-bosh-stemcell-${STEMCELL_VERSION}-google-kvm-ubuntu-xenial-go_agent.tgz" \
    -O stemcell.tgz

  om upload-stemcell -s stemcell.tgz

  rm stemcell.tgz

  om upload-product --product "tile-${TILE_VERSION}.pivotal"
  om stage-product --product-name=hello --product-version="${TILE_VERSION}"
  om apply-changes
cd - || exit 1