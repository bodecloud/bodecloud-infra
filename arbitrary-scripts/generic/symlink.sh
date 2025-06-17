#!/bin/bash

SYMLINK_PATH="/var/lib/containerd"
DEST_PATH="/mnt/blockvolume/k3s/containerd"

sudo rm -rvf $SYMLINK_PATH
sudo ln -s $DEST_PATH $SYMLINK_PATH 