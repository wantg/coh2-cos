cd ./src
rm -rf ../bin
mkdir -p ../bin/logs 
go generate 
go build -o ../bin/main
ln -fs ../assets ../bin
../bin/main -c ../config/config.yml
