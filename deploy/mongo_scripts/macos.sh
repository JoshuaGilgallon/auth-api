#!/bin/bash
# Check if MongoDB is running
if pgrep -x "mongod" > /dev/null
then
    echo "MongoDB is already running."
else
    echo "Starting MongoDB..."
    brew services start mongodb/brew/mongodb-community
fi
