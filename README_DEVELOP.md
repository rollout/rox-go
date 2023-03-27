# Developer Guide

## Run tests

### Unit tests
```
cd v5 && go test ./core/...
```

## Release a new v5 version

1. Update the semver version string returned by `LibVersion()` in `v5/core/client/device_properties.go`.
2. Make a PR, get it code reviewed and merged.
3. Update `master` and tag it with a new version number (the same one from step 1.). E.g.:
   ```shell
   git tag -a v5.0.5 -m "Release v5.0.5"
   ```
4. Push the tag to GitHub:
   ```shell
    git push origin v5.0.5
    ```
5. Profit.
