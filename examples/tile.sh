export TILE_VERSION="0.1.5"



cd examples/tile || exit 1
  kiln bake --version="${TILE_VERSION}"

  om upload-product --product tile-*.pivotal
  rm tile-*.pivotal
  om stage-product --product-name=hello --product-version="${TILE_VERSION}"
  om apply-changes
cd - || exit 1