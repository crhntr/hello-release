# Hello

This is a simple bosh release.

See examples/deployment for an example deployment manifest.

## Examples

### BOSH Manifest

See: ./examples/deployment/manifest.yml

- Configure BOSH (for example with `eval "$(bbl print-env)`
- Upload a release `bosh upload-release`
- Upload a stemcell  `bosh upload-stemcell "${RECENT_STEMCELL_URL}"
- Deploy `bosh deploy ./examples/deployment/manifest.yml`

### Tile

See: https://github.com/crhntr/hello-tile

- Create a release tarball and copy it to `../hello-tile/releases`
- (optional) Change the stemcell version in the Kilnfile.lock.
- Use `kiln bake --version=0.2.0` to create it.
- Upload to Ops Manger
- Configure the network stuff
- Hit apply changes... take a walk
