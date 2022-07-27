set -xe

cd tests/fakevim
make
cd -
mkdir -p ./tests/bin
mv tests/fakevim/bin/fakevim tests/bin/.
