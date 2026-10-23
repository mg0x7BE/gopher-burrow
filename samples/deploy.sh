#!/usr/bin/env bash
set -euo pipefail

readonly APP="inventory"
VERSION="${1:-$(git describe --tags --always)}"
TARGETS=(linux/amd64 linux/arm64 darwin/arm64)

log() {
    printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*"
}

for target in "${TARGETS[@]}"; do
    os="${target%/*}"
    arch="${target#*/}"
    log "building $APP $VERSION for $os/$arch"
    GOOS="$os" GOARCH="$arch" go build -o "dist/${APP}-${os}-${arch}" ./cmd/inventory
done

if [[ -n "${DEPLOY_HOST:-}" ]]; then
    scp dist/* "deploy@${DEPLOY_HOST}:/opt/${APP}/" && log "uploaded"
else
    echo "DEPLOY_HOST not set, skipping upload" >&2
fi

exit 0
