set -xe

cd tests/fakevim
make
cd -
mv tests/fakevim/bin/fakevim tests/bin/.
