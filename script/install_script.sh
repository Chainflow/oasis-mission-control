#!/bin/bash

set -e

cd $HOME

teleFalg="$1"
teleFlagValue="--remote-hosted"

echo "----------- Installing grafana -----------"

sudo -S apt-get install -y adduser libfontconfig1

wget https://dl.grafana.com/oss/release/grafana_6.7.2_amd64.deb

sudo -S dpkg -i grafana_6.7.2_amd64.deb

echo "------ Starting grafana server using systemd --------"

sudo -S systemctl daemon-reload

sudo -S systemctl start grafana-server

cd $HOME

echo "----------- Installing Influx -----------"

wget -qO- https://repos.influxdata.com/influxdb.key | sudo apt-key add -
source /etc/lsb-release
echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable" | sudo tee /etc/apt/sources.list.d/influxdb.list

sudo -S apt-get update && sudo apt-get install influxdb
sudo -S service influxdb start

cd $HOME

if [ "$teleFalg" != "$teleFlagValue" ];
then 
	echo "----------- Installing telegraf -----------------"
	
	sudo -S apt-get update && sudo apt-get install telegraf
	sudo -S service telegraf start

else
	echo "------remote-hosted enabled, so not downloading the telegraf--------"
fi

echo "------------Creating databases oasis and telegraf-------------"

curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE oasis"


echo "--------- Cloning oasis-validator-mission-control -----------"

cd go/src/github.com

git clone https://github.com/Chainflow/oasis-mission-control.git

cd oasis-mission-control

cp example.config.toml config.toml 

echo "------ Building and running the code --------"

go build && ./oasis-node