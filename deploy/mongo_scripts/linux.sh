#!/bin/bash
# Check if MongoDB is running
if systemctl is-active --quiet mongod
then
    echo "MongoDB is already running."
else
    echo "Starting MongoDB..."
    sudo systemctl start mongod
fi
