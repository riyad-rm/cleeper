# Cleeper

Cleeper is a tool to automatically shutdown resources within AWS.  
It currently supports:  
* EC2 instances : Stops the instances  
* RDS databases : Stops RDS instances and clusters  
* ASG autoscalling groups : Terminates instances within the ASG and suspends the Launch process  


## How to deploy

Make sure you have terraform, go and zip installed.  

Go to the deploy/ folder and run the deploy.sh script.   

This will compile the Lambda and run terraform to deploy it.

## How to use  

By default the lambda is deployed without triggers.  
You can either use a scheduled trigger to run it, invoke it from the cli or both.  

One way to run your account is to add two scheduled triggers, the first to start resources in the morning and the second to stop them every evening.  
A better way would be to stop the resources everyday and only start the resources you want when you need them by invoking the lambda from the cli, thus ensuring that forgotten resources are stopped and only started when needed.  

### Parameters

The Lambda accepts a number of parameters that can be used to customize its behavior. They are listed below.

1. action (required): this parameters controls the action taken by the lambda and can be one the 3 following values:
	1. start : will start the resources  
	2. stop: will stop the resources  
	3. list: a dry run which will list the resources that will stopped and started with the current parameters without actually stopping or starting them
2. regions (optional): a comma separated list of regions on which to act, ex: eu-west-1,eu-west-2. If not provided it will loop through all the AWS regions available. 
3. taggedOnly (optional): default value "true", can be "true" or "false", tells the lambda wether to ignore the tags or not when selecting which resources to stop/start
4. tagKeys (optional): default value "cleeper" , a comma separated list of tag keys to select resources to act on
5. tagValues (optional): default value "true", a comma separated list of tag values to select resources to act on

It is recommended to provide the list of regions you are using to avoid looping on empty regions and save on runtime cost.  

### Examples

To list the impacted resources in two regions:  
```json
{
  "action": "list",
  "regions": "eu-west-1,eu-west-2"
}
```  


If you want to shutdown all the resources that have the tag key **application** with either the value **app1** or **secondapp** in *eu-west-1*, you could pass the following parameters:  
```json
{
  "action": "stop",
  "regions": "eu-west-1",
  "tagKeys":"application",
  "tagValues":"app1,secondapp"
}
```  


If you want to shutdown all the resources regardless of their tag keys or tag value, simply set the taggedOnly parameter to false:  
```json
{
  "action": "stop",
  "taggedOnly": "false"
}
```  

The tags selection mechanism allows you to have precise control over which resources to stop or start, using the tags already in place.  

If you were to run a scheduled trigger you could use the json parameters above.  
In a cli scenario the lambda Invoke would look like this:  
```bash
aws lambda invoke --function-name cleeper --cli-binary-format raw-in-base64-out --payload '{"action":"list", "regions":"eu-west-1", "tagKeys":"cleeper", "tagValues":"val2,val1"}' --log-type Tail output | jq .LogResult -r | base64 -d
START RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889 Version: $LATEST
Working on region:  eu-west-1
ASGs to suspend
terraform-20250226115254216000000003
EC2 to terminate: 
i-0c4db9d0889e307fb
EC2 to stop:
i-06e32396035cf6c30
RDS Clusters to stop: 
aurora-cluster-demo
aurora-postgres-cluster-demo
RDS instances to stop: 
terraform-20250226124539700800000001
ASGs to resume
EC2 to start
RDS Clusters to start: 
RDS instances to start: 
END RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889
REPORT RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889	Duration: 592.21 ms	Billed Duration: 593 ms	Memory Size: 128 MB	Max Memory Used: 35 MB	
```  
In this example we are running a list action to see which resources would be stopped or started in the eu-west-1 region using custom tags.