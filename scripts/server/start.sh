#!/bin/bash

file_dir=$(dirname "$0")

# Ensure the script is running from the root of the server directory
echo "Moving to the directory of this script..."
cd $file_dir

echo "Moving to the root directory..."
cd ../../

# Start the server
echo "Starting the server..."
go run server/main.go