#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PWD=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
export GOPATH=$PWD:$GOPATH
##############################
######Install Dependence######
echo "Installing Dependence"
#go get github.com/gographics/imagick/imagick
##############################
#########Running Test#########
echo "Running Test"
# pkgs="\
# "
pkgs="\
 org.cny.sr/impl\
 org.cny.sr/mr\
"
echo "mode: set" > a.out
for p in $pkgs;
do
 echo $p
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $PWD/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html

echo "Build main"
go build -o sr org.cny.sr/main
