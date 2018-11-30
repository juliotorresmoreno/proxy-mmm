#!/bin/sh

go build -o /tmp/proxy-keos
sudo adduser proxy-keos --disabled-login --disabled-password --home /var/lib/proxy-keos --gecos ""
sudo cp proxy-keos.service /lib/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable proxy-keos.service
sudo systemctl start proxy-keos.service
sudo mv /tmp/proxy-keos /usr/bin/proxy-keos
sudo mkdir -p /etc/proxy-keos
sudo cp config.sample.json /etc/proxy-keos

