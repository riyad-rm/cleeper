package awsHandler

import (
 "github.com/aws/aws-sdk-go/aws"
 "github.com/aws/aws-sdk-go/service/autoscaling"
 "fmt"
)

// Be aware if you have replace unhealthy disabled
// Suspend ASGs
func suspendAsgs(client *autoscaling.AutoScaling, asgs []*autoscaling.Group){
	input := &autoscaling.ScalingProcessQuery{ScalingProcesses: []*string{aws.String("Launch")}}
	for _, asg := range asgs {
		fmt.Println("Suspending launch for : "+*asg.AutoScalingGroupName)
		input.SetAutoScalingGroupName(*asg.AutoScalingGroupName)
		client.SuspendProcesses(input)
	}
}

func resumeAsgs(client *autoscaling.AutoScaling, asgs []*autoscaling.Group){
	input := &autoscaling.ScalingProcessQuery{ScalingProcesses: []*string{aws.String("Launch")}}
	for _, asg := range asgs {
		fmt.Println("Resuming launch for : "+*asg.AutoScalingGroupName)
		input.SetAutoScalingGroupName(*asg.AutoScalingGroupName)
		client.ResumeProcesses(input)
	}

}



func evaluateASG(asg *autoscaling.Group, tag_keys *[]string, tag_values *[]string, suspend bool, tagged_only bool) bool{
	launchProcessSuspended := false
	for _, process := range asg.SuspendedProcesses {
			if *process.ProcessName == "Launch"{
				launchProcessSuspended = true
				break
			}
	}
	if suspend && launchProcessSuspended {
		// no need to suspend if already suspended
		return false
	}
	if !suspend && !launchProcessSuspended {
		// no need to resume if not suspended
		return false
	}
	if !tagged_only{
	 	return true
	 }
	for _, tag :=range asg.Tags {
		if stringInList(tag_keys, *tag.Key){
			if stringInList(tag_values, *tag.Value){
				return true
			}
		}
	}
	return false
}

func listASGInstances(asg *autoscaling.Group) []*autoscaling.Instance{
	return asg.Instances
}

// Launching this two times in a row the terminated instance is still part of the ASG and is returned
// Terminating an already terminated instance does not result in errors
func listASGsInstancesIds(asgs []*autoscaling.Group) []*string{
	var instances []*string
	for _, asg := range asgs{
		for _, instance := range asg.Instances{
			instances = append(instances, instance.InstanceId)
		}
	}
	return instances
}

// List of ASG filtered based on tag
// ignore already suspended or non suspended ASGs
func listFilteredASG(client *autoscaling.AutoScaling, tag_keys *[]string, tag_values *[]string, tagged_only bool, suspend bool) []*autoscaling.Group{
	//ListAllASG(asgSvc)
	input := &autoscaling.DescribeAutoScalingGroupsInput{MaxRecords : aws.Int64(MAX_ASG_INSTANCES)}
	var asgList []*autoscaling.Group
	for true {
		result, err := client.DescribeAutoScalingGroups(input)
		if err != nil{
			fmt.Println(err)
		}

		for _, asg := range result.AutoScalingGroups{
			if evaluateASG(asg, tag_keys, tag_values, suspend, tagged_only){
				asgList = append(asgList, asg)
			}
		}
		if result.NextToken != nil {
			input.SetNextToken(*result.NextToken)
		} else {
			break
		}
	}
	//suspendAsgs(client, asgList)
	return asgList
}