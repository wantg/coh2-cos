cd ./src
mkdir -p ../bin/logs 
go generate 
go build -o ../bin/main
rm -rf ../bin/assets
ln -fs ../assets ../bin
../bin/main -c ../config/config.yml
