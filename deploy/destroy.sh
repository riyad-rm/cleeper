#!/bin/bash

cd terraform
terraform init
terraform destroy --auto-approve
rm ../../build/bin/bootstrap
rm ../../build/bin/cleeper.zip