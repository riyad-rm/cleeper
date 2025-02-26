package awsHandler

import (
 "github.com/aws/aws-sdk-go/aws/session"
 "github.com/aws/aws-sdk-go/service/ec2"
 "github.com/aws/aws-sdk-go/aws"
 "github.com/aws/aws-sdk-go/service/autoscaling"
 "github.com/aws/aws-sdk-go/service/rds"
 "fmt"
 "strings"
)


// General
// Test if list of resources is empty before any call to a stop/start function
// Sending an empty list to these function will trigger a warning err

// Region handlers
// ! need to return proper error
func StopRegion(sess *session.Session, tag_keys *[]string, tag_values *[]string, region string, dry_run bool, tagged_only bool){
	/*
	*		ASG
	*/
	// List ASG
	asgSvc := autoscaling.New(sess, aws.NewConfig().WithRegion(region))
	asgList := listFilteredASG(asgSvc, tag_keys, tag_values, tagged_only, true)
	ec2InstancesToTerminate := listASGsInstancesIds(asgList)
	fmt.Println("ASGs to suspend")
	for _, str := range asgList {
		fmt.Println(*str.AutoScalingGroupName)
	}
		if !dry_run {
		suspendAsgs(asgSvc, asgList)
	}
	// we terminate the ASG instances in the following bloc dealing with EC2


	/*
	*		EC2
	*/
	ec2Svc := ec2.New(sess, aws.NewConfig().WithRegion(region))
	// Running instances only
	filters := []*ec2.Filter{
	            &ec2.Filter{
	                Name: aws.String("instance-state-name"),
	                Values: []*string{
	                    aws.String("running"),
	                },
	            },
	        }
	ec2List := listFilteredEC2Instances(ec2Svc, tag_keys, tag_values, tagged_only, filters)
	// !!! This prints a weired id needs checking 
	fmt.Println("EC2 to terminate: ")
	for _, str := range ec2InstancesToTerminate {
		fmt.Println(*str)
	}
	if len(ec2InstancesToTerminate) > 0 && !dry_run {
		terminateEC2output, err := terminateEC2Instances(ec2Svc, ec2InstancesToTerminate)
		if err!=nil {
			fmt.Println(err)
			fmt.Println(terminateEC2output)
		}
	}
	fmt.Println("EC2 to stop:")
	for _, str := range ec2List {
		fmt.Println(*str)
	}
	if len(ec2List) > 0  && !dry_run{
		stopEC2output, err := stopEC2Instances(ec2Svc, ec2List)
		if err!=nil {
			fmt.Println(err)
			fmt.Println(stopEC2output)
		}
	}

	/*
	*		RDS
	*/
	rdsSvc := rds.New(session.New(), aws.NewConfig().WithRegion(region))
	// list rds clusters
	rdsClustersList := listFilteredRDSClusters(rdsSvc, tag_keys, tag_values, tagged_only, true)
	fmt.Println("RDS Clusters to stop: ")
	for _, instance := range rdsClustersList {
		fmt.Println(*instance)
	}
	if !dry_run {
		stopRDSClusters(rdsSvc, rdsClustersList)
	}
	// rds instances
	rdsInstanceList := listFilteredRDSInstances(rdsSvc, tag_keys, tag_values, tagged_only, true)
	fmt.Println("RDS instances to stop: ")
	for _, instance := range rdsInstanceList {
		fmt.Println(*instance)
	}
	if !dry_run {
		stopRDSInstances(rdsSvc, rdsInstanceList)
	}
}

func StartRegion(sess *session.Session, tag_keys *[]string, tag_values *[]string, region string, dry_run bool, tagged_only bool){

	/*
	*		ASG
	*/
	asgSvc := autoscaling.New(sess, aws.NewConfig().WithRegion(region))
	asgList := listFilteredASG(asgSvc, tag_keys, tag_values, tagged_only,  false)
	fmt.Println("ASGs to resume")
	for _, str := range asgList {
		fmt.Println(*str.AutoScalingGroupName)
	}
	if len(asgList)>0 && !dry_run {
		resumeAsgs(asgSvc, asgList)
	}

	/*
	*		EC2
	*/
	ec2Svc := ec2.New(sess, aws.NewConfig().WithRegion(region))
	// stopped instances only
	filters := []*ec2.Filter{
	            &ec2.Filter{
	                Name: aws.String("instance-state-name"),
	                Values: []*string{
	                    aws.String("stopped"),
	                },
	            },
	        }
	ec2List := listFilteredEC2Instances(ec2Svc, tag_keys, tag_values, tagged_only, filters)
	fmt.Println("EC2 to start")
	for _, str := range ec2List {
		fmt.Println(*str)
	}
	if len(ec2List)>0  && !dry_run{
		startEC2output, err := startEC2Instances(ec2Svc, ec2List)
		if err!=nil {
			fmt.Println(err)
			fmt.Println(startEC2output)
		}
	}

	/*
	*		RDS
	*/
	rdsSvc := rds.New(session.New(), aws.NewConfig().WithRegion(region))
	// list rds clusters
	rdsClustersList := listFilteredRDSClusters(rdsSvc, tag_keys, tag_values, tagged_only, false)
	fmt.Println("RDS Clusters to start: ")
	for _, instance := range rdsClustersList {
		fmt.Println(*instance)
	}
	if !dry_run {
		startRDSClusters(rdsSvc, rdsClustersList)
	}
	
	// rds instances
	rdsInstanceList := listFilteredRDSInstances(rdsSvc, tag_keys, tag_values, tagged_only, false)
	fmt.Println("RDS instances to start: ")
	for _, instance := range rdsInstanceList {
		fmt.Println(*instance)
	}
	if !dry_run {
		startRDSInstances(rdsSvc, rdsInstanceList)
	}
}


func ListRegion(sess *session.Session, tag_keys *[]string, tag_values *[]string, region string, tagged_only bool){
	// listing is always a dry_run
	// Stop
	StopRegion(sess, tag_keys, tag_values, region, true, tagged_only)
	// Start
	StartRegion(sess, tag_keys, tag_values, region, true, tagged_only)
}

// Global handlers
// ! need to return a proper status on error
func Action(action LambdaTrigger){
	// load tags config
	tag_keys := []string{"cleeper"}
	tag_values := []string{"true"}
	if action.TagKeys != ""{
		tag_keys = strings.Split(action.TagKeys, ",")
	}
	if action.TagValues != ""{
		tag_values = strings.Split(action.TagValues, ",")
	}

	// session to be used by all clients
	sess := session.Must(session.NewSessionWithOptions(session.Options{
	    SharedConfigState: session.SharedConfigEnable,
	}))

	// Load regions config
	regions := []string{}
	if action.Regions!="" {
		regions = strings.Split(action.Regions, ",")
	} else {
		regions, _ = ListAllRegions(sess)
	}
	tagged_only := true
	if action.TaggedOnly == "false" {
		tagged_only = false
	}

	// walk through the regions
	for _, region := range regions {
		fmt.Println("Working on region: ", region)
		switch action.Action {
			case "list": 
				ListRegion(sess, &tag_keys, &tag_values, region, tagged_only)
			case "stop": 
				StopRegion(sess, &tag_keys, &tag_values, region, false, tagged_only)
			case "start": 
				StartRegion(sess, &tag_keys, &tag_values, region, false, tagged_only)
			default:
				fmt.Println("Unsupported action:",action.Action)
		}
		
	}
}

func FullStop(){

}

func initConfig(){
	
}