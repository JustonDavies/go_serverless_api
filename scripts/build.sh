#!/usr/bin/env bash
#-- Shared variables ----------
project_owner="justondavies"
project_name="go_serverless_api"

process="build"

build_directory="build"

#-- Begin ----------
echo "Preparing ${process} shell for ${project_owner}/${project_name}"

#-- Create image with build tools ----------
image_tag="${project_owner}/${project_name}:build"
working_directory="/${process}/${project_owner}/${project_name}"
sudo docker build                                                          \
  --network host                                                           \
  --file dockerfiles/${process}.docker                                     \
  --tag ${image_tag}                                                       \
  --build-arg WORKING_DIRECTORY=${working_directory}                       \
  ./

#-- Create container for build / extract ----------
container_name="${project_owner}_${project_name}_ephemeral_${process}_context"
sudo docker run                                                            \
  --name ${container_name}                                                 \
  --volume $PWD/${build_directory}:${working_directory}/${build_directory} \
  --interactive                                                            \
  --tty                                                                    \
  --rm                                                                     \
  ${image_tag}

#-- Clean up ----------
sudo chown -R ${USER}:${USER} ./${build_directory}
chmod -R 777 ./${build_directory}
