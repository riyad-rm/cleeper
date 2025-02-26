package awsHandler

import (
 "github.com/aws/aws-sdk-go/aws/session"
 "github.com/aws/aws-sdk-go/service/ec2"
 "github.com/aws/aws-sdk-go/aws"
 "fmt"
)

/*
Convention:
	Typo:
		[] = mandatory group
		() = optional group
		<> = to be replaced by a value
	Functions Naming:
		List[All|<Relevant criteria>]<ResourceType>(<AdditonalInfo>)
		Ex: ListAllEC2Instances, ListAllRDSInstances, ListAllRDSClusters ...
	Input:
		client = an AWS client for the desired resource type
		desiredProperties = list of the resources properties to return ex: id, arn : name , state
	Return:
		return references not variable copies

Guidelines:
	Sessions must be shared between clients wherever possible (ec2, s3 ...)
	Clients must be shared wherever possible (ec2: instances, asg, loadbalancers ...)
*/


// TODO list All resources

func listAllEC2Instances(){
}

func evaluateEC2(instance *ec2.Instance, tag_keys *[]string, tag_values *[]string, tagged_only bool) bool{
	haveTag := false
	for _, tag :=range instance.Tags {
		// reject instances that are part of an ASG as they will be terminated
		if *tag.Key == "aws:autoscaling:groupName" {
			return false
		}
		if stringInList(tag_keys, *tag.Key){
			if stringInList(tag_values, *tag.Value){
				haveTag = true
				// do not return here to make sure we have seen all tags and none of them is aws:autoscaling:groupName
			}
		}
	}
	if !tagged_only {
		// If we don't care about the tag, we accept all instances
		return true
	}
	return haveTag
}

func listFilteredEC2Instances(client *ec2.EC2, tag_keys *[]string, tag_values *[]string, tagged_only bool, filters []*ec2.Filter) []*string{
	var ec2List []*string
	input := &ec2.DescribeInstancesInput{MaxResults : aws.Int64(MAX_EC2_INSTANCES)}
	input.SetFilters(filters)
	for true {
		result, err := client.DescribeInstances(input)
		if err != nil {
			fmt.Println(err)
		}
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				if evaluateEC2(instance, tag_keys, tag_values, tagged_only){
					ec2List = append(ec2List, instance.InstanceId)
				}
			}
		}
		if result.NextToken != nil {
			input.SetNextToken(*result.NextToken)
		} else {
			break
		}
	}
	return ec2List
}

func ListAllRegions(sess *session.Session) ([]string, error){
	input := &ec2.DescribeRegionsInput{}
	input.SetAllRegions(false)
	client := ec2.New(sess)
	output, err := client.DescribeRegions(input)
	regions := make([]string, len(output.Regions))
	for i, region := range output.Regions{
		regions[i] = *region.RegionName
	}
	return regions, err
}


func terminateEC2Instances(client *ec2.EC2, instancesIds []*string) (*ec2.TerminateInstancesOutput, error){
	input := &ec2.TerminateInstancesInput{}
	input.SetInstanceIds(instancesIds)
	return client.TerminateInstances(input)
}

func stopEC2Instances(client *ec2.EC2, instancesIds []*string) (*ec2.StopInstancesOutput, error){
	input := &ec2.StopInstancesInput{}
	input.SetInstanceIds(instancesIds)
	return client.StopInstances(input)
}

func startEC2Instances(client *ec2.EC2, instancesIds []*string) (*ec2.StartInstancesOutput, error){
	input := &ec2.StartInstancesInput{}
	input.SetInstanceIds(instancesIds)
	return client.StartInstances(input)
}