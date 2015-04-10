#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PWD=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
export GOPATH=$PWD:$GOPATH
o_dir=build
rm -rf $o_dir
mkdir $o_dir

#### Package fvm ####
v_srv=0.0.1
o_srv=$o_dir/sr
mkdir $o_srv
mkdir $o_srv/conf
mkdir $o_srv/www
go build -o $o_srv/sr org.cny.sr/main
cp srid $o_srv
cp conf/sr.properties $o_srv/conf
if [ "$1" != "" ];then
	curl -o $o_srv/srvd_i $1
	chmod +x $o_srv/srvd_i
	echo "./srvd_i \$1 srid" >$o_srv/install.sh
	chmod +x $o_srv/install.sh
fi 
cd $o_dir
zip -r sr.zip sr
cd ../
echo "Package sr..."