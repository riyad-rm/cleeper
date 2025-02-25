#!/bin/bash

cd ../build
./build.sh
cd ../deploy/terraform
terraform init
terraform apply --auto-approve