set -ex

WIKICMD_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

export GH_TOKEN=${GH_KEY}

echo $WIKICMD_VERSION

gh release create "$WIKICMD_VERSION" -t "$WIKICMD_VERSION"

RELEASE_FILES=$(ls ./release_downloads/*.zip)

for RELEASE_FILE in $RELEASE_FILES
do
    gh release upload "$WIKICMD_VERSION" "$RELEASE_FILE"
done
