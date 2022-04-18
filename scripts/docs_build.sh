set -ex

URL_MDBOOK="https://github.com/rust-lang/mdBook/releases/download/v0.4.18/mdbook-v0.4.18-x86_64-unknown-linux-gnu.tar.gz"

if [[ -z "$WIKICMD_DOCS_BUILD_SAVE_AS" ]]
then
    WIKICMD_DOCS_BUILD_SAVE_AS="$(pwd)""/docs.tgz"
fi

DIR_TMP=$(mktemp -d)

cp -r docs "$DIR_TMP""/."

cd "$DIR_TMP"

wget -O mdbook.gz "$URL_MDBOOK"
tar xzvf ./mdbook.gz
cd docs
../mdbook build

tar czvf "$WIKICMD_DOCS_BUILD_SAVE_AS" -C ./book .