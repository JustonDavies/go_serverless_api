#!/usr/bin/env bash
#TODO: This doesn't quite work right now, come back later
#echo Building justondavies/go_serverless_api:deploy
#
#sudo docker build                                                          \
#  --network host                                                           \
#  --file dockerfiles/deploy.docker                                         \
#  --tag justondavies/go_serverless_api:deploy                              \
#  ./
#
#sudo docker create                                                         \
#  --name deploy_state_extract                                              \
#  justondavies/go_serverless_api:deploy
#
#rm -rf ./build/*
#
#sudo docker cp                                                             \
#  build_extract:/go/src/github.com/justondavies/go_serverless_api/build    \
#  ./
#
#sudo docker rm -f build_extract
#
#sudo chown -R $USER:$USER ./build
#chmod -R 777 ./build
