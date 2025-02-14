# Cleeper

Cleeper is a tool to automatically shutdown ressources within aws.  
It currently supports:  
* EC2 instances : Stops the instances  
* RDS databases : Stops RDS instances and clusters  
* ASG autoscalling groups : Terminates instances within the ASG and suspends the Launch process  


## How to deploy

Go to deploy/terraform and run de deploy.sh script.  
Make sure you have golang and zip installed before.  

  
This will compile the Lambda an run terraform to deploy it.

## How to use  