set -ex

URL_MDBOOK="https://github.com/rust-lang/mdBook/releases/download/v0.4.49/mdbook-v0.4.49-x86_64-unknown-linux-gnu.tar.gz"
LATEST_VERSION=$(git describe --tags --abbrev=0 origin/master)

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

cd book
FILES_TO_REPLACE=$(grep -rl '%WIKICMD_VERSION%' . | grep '\.html$')
for FILE in "$FILES_TO_REPLACE"
do
    sed -i "s/%WIKICMD_VERSION%/$LATEST_VERSION/g" $FILE
done
cd ..

tar czvf "$WIKICMD_DOCS_BUILD_SAVE_AS" -C ./book .
