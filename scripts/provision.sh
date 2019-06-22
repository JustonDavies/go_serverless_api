#!/usr/bin/env bash
#-- Shared variables ----------
project_owner="justondavies"
project_name="go_serverless_api"

process="provision"

secrets_directory="configs/secrets"

state_directory=".state"
tool_state_directory="terraform.tfstate.d"

#-- Begin ----------
echo "Preparing ${process} shell for ${project_owner}/${project_name}"

#-- Backup state files --------
version=`date +%d_%m_%y_%H_%M_%S`
mkdir -p ${state_directory}/backup/${version}_${process}
cp -r ${state_directory}/${tool_state_directory} ${state_directory}/backup/${version}_${process}/${tool_state_directory}

#-- Create image with process tools ----------
image_tag="${project_owner}/${project_name}:${process}_context"
working_directory="/${process}/${project_owner}/${project_name}"
sudo docker build                                                                           \
  --network host                                                                            \
  --file dockerfiles/${process}.docker                                                      \
  --tag ${image_tag}                                                                        \
  --build-arg WORKING_DIRECTORY=${working_directory}                                        \
  ./

#-- Create container with process context ----------
container_name="${project_owner}_${project_name}_ephemeral_${process}_context"
sudo docker run                                                                             \
  --name ${container_name}                                                                  \
  --volume $PWD/${state_directory}/${tool_state_directory}:${working_directory}/${tool_state_directory} \
  --volume $PWD/${secrets_directory}:${working_directory}/${secrets_directory}              \
  --volume $HOME/.aws:/root/.aws                                                            \
  --interactive                                                                             \
  --tty                                                                                     \
  --rm                                                                                      \
  ${image_tag}

#-- Clean up ----------
sudo chown -R ${USER}:${USER} ./${state_directory}
