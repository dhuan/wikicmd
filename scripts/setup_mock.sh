set -xe

MOCK_DOWNLOAD_URL="https://github.com/dhuan/mock/releases/download/v0.0.2/mock_v0.0.2_linux-386.tar.gz"

TMP_PATH=$(mktemp -d)
wget -P "$TMP_PATH" "$MOCK_DOWNLOAD_URL"
tar xzvf "$TMP_PATH"/*.tar.gz -C "$TMP_PATH" 
mkdir -p ./tests/bin
mv "$TMP_PATH""/mock" ./tests/bin/.
