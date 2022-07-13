#!/bin/sh

apt-get -y update

# Install wireguard
apt install -y wireguard

# install netclient
curl -sL 'https://apt.netmaker.org/gpg.key' | sudo tee /etc/apt/trusted.gpg.d/netclient.asc
curl -sL 'https://apt.netmaker.org/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/netclient.list
sudo apt update -y
sudo apt install -y netclient

systemctl enable netclient

# Join the node to the mesh network
netclient join -t eyJhcGljb25uc3RyaW5nIjoiYXBpLm5tLjU0LTIzNy0zOC0xMzUubmlwLmlvOjQ0MyIsIm5ldHdvcmsiOiJrM3MiLCJrZXkiOiJmZDk2ZTUwMGVmMzA5NDhjIiwibG9jYWxyYW5nZSI6IiJ9

# Install K3S and connect to the control plane
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.21.7+k3s1 K3S_URL=https://172.30.16.1:6443 K3S_TOKEN=K10ece84e5bc0ffa0759311c80e6eaba2db7e486a745bd02f133f01838b3dd5a3ea::server:823348c08eeeafa03058f0fe9294ef9c sh -

# Configure Flannel CNI

export IP=$(ip addr show nm-k3s | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')

DIR="/etc/systemd/system/k3s.service.d"
FILE="/etc/systemd/system/k3s.service.d/network.conf"
if [ ! -d "$DIR" ]; then
mkdir /etc/systemd/system/k3s.service.d
fi

if [ ! -f "$FILE" ]; then
touch /etc/systemd/system/k3s.service.d/network.conf
fi

echo `cat <<EOF
[Service] \n
ExecStart= \n
ExecStart=/usr/local/bin/k3s agent --node-ip $IP --flannel-iface=nm-k3s
EOF` > /etc/systemd/system/k3s.service.d/network.conf
systemctl daemon-reload
systemctl restart k3s-agent.service
