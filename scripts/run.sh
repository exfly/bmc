#!/usr/bin/env bash

set -o errtrace
set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# cd "$(dirname "$0")"

# sudo apt install jq

mkdir -p tmp/rootfs/host/
mkdir -p tmp/host/

tar -C tmp/ -xvf /vagrant/tmp/runtime.tar
ROOTFS_SHA256=$(cat tmp/repositories | jq '."exfly/skopeo".dev' -r)
# tar -C tmp/ -xvf /vagrant/tmp/alpine-edge.tar
# ROOTFS_SHA256=$(cat tmp/repositories | jq '.alpine.edge' -r)


tar -C tmp/rootfs/ -xvf tmp/${ROOTFS_SHA256}/layer.tar # 获得 rootfs

cp /vagrant/config.json tmp/

cp /etc/resolv.conf tmp/rootfs/etc/resolv.conf
cp /etc/hosts tmp/rootfs/etc/hosts

chmod +x /vagrant/bin/runc.amd64
cd tmp
/vagrant/bin/runc.amd64 run --no-pivot mycontainerid
