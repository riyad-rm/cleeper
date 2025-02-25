# Cleeper

Cleeper is a tool to automatically shutdown ressources within aws.  
It currently supports:  
* EC2 instances : Stops the instances  
* RDS databases : Stops RDS instances and clusters  
* ASG autoscalling groups : Terminates instances within the ASG and suspends the Launch process  


## How to deploy

Make sure you have terraform, go and zip installed.  

Go to the deploy/ folder and run the deploy.sh script.   

This will compile the Lambda and run terraform to deploy it.

## How to use  

It is a best practice to set the regions you want the lambda to work on to avoid looping on empty regions  