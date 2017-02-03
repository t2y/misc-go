#!/bin/bash

if [ $# -ne 2 ]; then
  echo "usage: init.sh your-host-fqdn mybucket"
  exit 1
fi

endpoint="$1"
bucket="$2"

# create bucket
./main -command putBucket -bucket "$bucket" -endpoint "$endpoint"

# get bucket after it created above
./main -command getBucket -bucket "$bucket" -endpoint "$endpoint"

