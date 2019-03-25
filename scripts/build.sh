#!/usr/bin/env bash
echo Building justondavies/go_serverless_api:build

sudo docker build                                                          \
  --network host                                                           \
  --file dockerfiles/build.docker                                          \
  --tag justondavies/go_serverless_api:build                               \
  ./

sudo docker create                                                         \
  --name build_extract                                                     \
  justondavies/go_serverless_api:build

rm -rf ./build/*

sudo docker cp                                                             \
  build_extract:/go/src/github.com/justondavies/go_serverless_api/build    \
  ./

sudo docker rm -f build_extract

sudo chown -R $USER:$USER ./build
chmod -R 777 ./build
