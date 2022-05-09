set -ex

WIKICMD_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

TARGETS=(
    "linux,386"
    "linux,amd64"
    "linux,arm"
    "linux,arm64"
    "darwin,amd64"
)

rm -rf ./release_downloads

mkdir ./release_downloads

for TARGET in "${TARGETS[@]}"
do
    GOOS=$(echo $TARGET | cut -d "," -f 1)
    GOARCH=$(echo $TARGET | cut -d "," -f 2)

    TARGET_NAME="${GOOS}-${GOARCH}"

    printf "Generating build for ${TARGET_NAME}\n"

    TARGET_PATH="./release_downloads/$TARGET_NAME"

    mkdir $TARGET_PATH

    cp ./README.md $TARGET_PATH/.
    cp ./LICENSE $TARGET_PATH/.

    TMP_BKP=$(mktemp)
    cp internal/wikicmd/wikicmd.go "$TMP_BKP"
    sed -i "s/__VERSION__/$WIKICMD_VERSION/g" internal/wikicmd/wikicmd.go
    sed -i "s/__GOOS__/$GOOS/g" internal/wikicmd/wikicmd.go
    sed -i "s/__GOARCH__/$GOARCH/g" internal/wikicmd/wikicmd.go

    GOOS=$GOOS GOARCH=$GOARCH go build -o $TARGET_PATH/wikicmd

    cp "$TMP_BKP" internal/wikicmd/wikicmd.go
done

TARGET_FOLDERS=$(ls ./release_downloads)

for TARGET_FOLDER in ${TARGET_FOLDERS[@]}
do
    zip "./release_downloads/wikicmd_${WIKICMD_VERSION}_${TARGET_FOLDER}.zip" -j ./release_downloads/"${TARGET_FOLDER}"/*
done
