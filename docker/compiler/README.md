Modified from https://github.com/bep/dockerfiles/tree/master/ci-goreleaser

When the dockerfile is changed, the version number should be incremented in the Makefile and the new version tag should be pushed to docker hub. The `scripts/cross-compile.sh` script should also be updated to use the new version number tag, and `.travis.yml` needs to be updated to pull the correct image tag.

A MacOS univeral binary can be created using `lipo -create -output stash-osx-universal stash-osx stash-osx-applesilicon`, available in the image.