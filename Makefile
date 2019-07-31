# This needs to match internal/diagbin.assetName.
ASSET_NAME := diagbin

binaries:
	go run ./internal/embed-binaries \
		-build-from ./lantern-diagnostics \
		-embed-to ./internal/diagbin \
		-pkg diagbin \
		-prefix embedded_command \
		-asset-name diagbin