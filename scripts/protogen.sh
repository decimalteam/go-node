#!/bin/bash

proto_dirs=$(find ./proto -name '*.proto' -print0 | xargs -0 -n1 dirname | uniq)

for dir in $proto_dirs; do
  buf protoc $(find "${dir}" -maxdepth 1 -name "*.proto")
done
