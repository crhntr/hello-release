# Hello

This is a simple bosh release.

See examples/deployment for an example deployment manifest.

## Examples

### BOSH Manifest

See: ./examples/deployment/manifest.yml

- Configure BOSH (for example with `eval "$(bbl print-env)`
- Upload a release `bosh upload-release`
- Upload a stemcell  `bosh upload-stemcell https://bosh.io/d/stemcells/bosh-google-kvm-ubuntu-xenial-go_agent?v=621.117`
- Deploy `bosh deploy ./examples/deployment/manifest.yml`

### Tile

See: ./examples/tile/.

- Create a release tarball and copy it to `./examples/tile/releases`
- (optional) Change the stemcell version in the Kilnfile.lock.
- Use `kiln bake` to create it.
- Upload to Ops Manger
- Configure the network stuff
- Hit apply changes... take a walk
