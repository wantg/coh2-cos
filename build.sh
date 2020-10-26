rm -rf dist
mkdir dist

cd hq/src
# go build -o ../dist/hq/hq .
go generate
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../../dist/hq/hq.exe .
cp -r ../assets ../../dist/hq
cp -r ../config/config.yml ../../dist/config.yml

cd ../../
# ./node_modules/.bin/electron-builder --config electron-builder.yml --win --mac
./node_modules/.bin/electron-builder --config electron-builder.yml --win
rm -rf dist/win-unpacked dist/mac dist/.icon-icns
