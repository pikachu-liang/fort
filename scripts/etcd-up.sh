#!/usr/bin/env bash
# A script that run a single node etcd cluster

# Use the host IP address when configuring etcd
export NODE1=192.168.1.21
# Configure a Docker volume to store etcd data
docker volume create --name etcd-data
export DATA_DIR="etcd-data"

REGISTRY=quay.io/coreos/etcd
# available from v3.2.5
REGISTRY=gcr.io/etcd-development/etcd
# pull image first
docker pull ${REGISTRY}:latest
# run a single node etcd
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380 \
  2>&1 &

# TODO: check etcd status in a for loop and timeout
echo -e "sleep 5s to wait docker starting etcd and etcd initialization"
sleep 5

# check if etcd is started successfully
member_list=$(docker exec etcd /bin/sh -c "export ETCDCTL_API=3 && /usr/local/bin/etcdctl member list")
if [[ $member_list == *" started, node1"* ]]; then
    echo "etcd started successfully. etcd member list:"
    echo "$member_list"
else
    echo "failed to start etcd."
    exit 1
fi