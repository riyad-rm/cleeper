package awsHandler

import (
 //"github.com/aws/aws-sdk-go/aws"
 "github.com/aws/aws-sdk-go/service/rds"
 "fmt"
 "strings"
)


// rds instances
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_StopInstance.html
func evaluateRDSInstance(instance *rds.DBInstance, tag_keys *[]string, tag_values *[]string, stop bool, tagged_only bool) bool{
	// Aurora DB can't be stoped or started
	if strings.HasPrefix(*instance.Engine, "aurora") {
		fmt.Println("Skip rds instance with aurora engine")
		return false
	}

	// ignore sql server multi az
	if strings.HasPrefix(*instance.Engine, "sqlserver") && *instance.MultiAZ == true {
		fmt.Println("Skip sqlserver engine with multi AZ deployment")
		return false
	}

	// ignore read replicas and replicated DBs
	if instance.ReadReplicaSourceDBInstanceIdentifier != nil {
		fmt.Println("Skip replica database ", *instance.DBInstanceIdentifier)
		return false
	}
	if len(instance.ReadReplicaDBInstanceIdentifiers) > 0 {
		fmt.Println("Skip Read replicated database ", *instance.DBInstanceIdentifier)
		return false
	}
	// Clusters are handled by antoher function
	if instance.DBClusterIdentifier != nil {
		fmt.Println("Skip instance that is part of a cluster")
		return false
	}

	// don't strop stopped instances and don't start available instances
	if stop {
		if ! stringInList(&[]string{"available"}, *instance.DBInstanceStatus){
			fmt.Println("Skip rds instance not available ",*instance.DBInstanceIdentifier)
			return false
		}
	} else {
		if ! stringInList(&[]string{"stopped"}, *instance.DBInstanceStatus){
				fmt.Println("Skip rds instance not stopped cannot be started ",*instance.DBInstanceIdentifier)
				return false
		}
	}

	// evaluate tags
	haveTag := false
	for _, tag := range instance.TagList {
		if stringInList(tag_keys, *tag.Key){
			if stringInList(tag_values, *tag.Value){
				haveTag = true
				break
			}
		}
	}
	if haveTag && tagged_only{
	 	return true
	 }
	haveTag = !haveTag
	return haveTag
}

func listFilteredRDSInstances(client *rds.RDS, tag_keys *[]string, tag_values *[]string, tagged_only bool, action bool) []*string{
	input:= &rds.DescribeDBInstancesInput{}
	input.SetMaxRecords(MAX_RDS_INSTANCES)
	
	var DBInstances []*string
	for true {
		result, err := client.DescribeDBInstances(input)
		if err != nil {
			fmt.Println(err)
		}
		for _, instance := range result.DBInstances {
			if evaluateRDSInstance(instance, tag_keys, tag_values, action, tagged_only){
				DBInstances = append(DBInstances, instance.DBInstanceIdentifier)
			}
		}
		if result.Marker != nil{
			input.SetMarker(*result.Marker)
		} else{
			break
		}
	}
	return DBInstances
}

func stopRDSInstances(client *rds.RDS, dbIdentifiers []*string){
	input := &rds.StopDBInstanceInput{}
	for _, db := range dbIdentifiers{
		input.SetDBInstanceIdentifier(*db)
		output, err := client.StopDBInstance(input)
		if err != nil {
			fmt.Println("Error while stoping instance ", *db)
			fmt.Println(err)
			fmt.Println(output)
		}
	}

}

func startRDSInstances(client *rds.RDS, dbIdentifiers []*string){
	input := &rds.StartDBInstanceInput{}
	for _, db := range dbIdentifiers{
		input.SetDBInstanceIdentifier(*db)
		output, err := client.StartDBInstance(input)
		if err != nil {
			fmt.Println("Error while starting instance ", *db)
			fmt.Println(err)
			fmt.Println(output)
		}
	}
}


// rds cluster

func evaluateRDSCluster(instance *rds.DBCluster, tag_keys *[]string, tag_values *[]string, stop bool, tagged_only bool) bool{
	//fmt.Println(*instance)
	if *instance.EngineMode != "provisioned" {
		fmt.Println("Skip non provisioned cluster")
		return false
	} 
	// don't stop stopped instances and don't start available instances
	if stop {
		if stringInList(&[]string{"stopped","stopping","deleting"}, *instance.Status){
			fmt.Println("Skip rds cluster not available ",*instance.DbClusterResourceId)
			return false
		}
	} else {
		if ! stringInList(&[]string{"stopped","stopping"}, *instance.Status){
			fmt.Println("Skip rds cluster not stopped ",*instance.DbClusterResourceId)
			return false
		}
	}

	// evaluate tags
	haveTag := false
	for _, tag := range instance.TagList {
		if stringInList(tag_keys, *tag.Key){
			if stringInList(tag_values, *tag.Value){
				haveTag = true
				break
			}
		}
	}
	if haveTag && tagged_only{
	 	return true
	 }
	haveTag = !haveTag
	return haveTag
}

func listFilteredRDSClusters(client *rds.RDS, tag_keys *[]string, tag_values *[]string, tagged_only bool, action bool) []*string{
	input:= &rds.DescribeDBClustersInput{}
	input.SetMaxRecords(MAX_RDS_INSTANCES)
	
	var DBClusters []*string
	for true {
		result, err := client.DescribeDBClusters(input)
		if err != nil {
			fmt.Println(err)
		}
		for _, instance := range result.DBClusters {
			if evaluateRDSCluster(instance, tag_keys, tag_values, action, tagged_only){
				DBClusters = append(DBClusters, instance.DBClusterIdentifier)
			}
		}
		if result.Marker != nil{
			input.SetMarker(*result.Marker)
		} else{
			break
		}
	}
	return DBClusters
}

func stopRDSClusters(client *rds.RDS, dbIdentifiers []*string){
	input := &rds.StopDBClusterInput{}
	for _, db := range dbIdentifiers{
		input.SetDBClusterIdentifier(*db)
		output, err := client.StopDBCluster(input)
		if err != nil {
			fmt.Println("Error while stopping cluster ", *db)
			fmt.Println(err)
			fmt.Println(output)
		}
	}
}
func startRDSClusters(client *rds.RDS, dbIdentifiers []*string){
	input := &rds.StartDBClusterInput{}
	for _, db := range dbIdentifiers{
		input.SetDBClusterIdentifier(*db)
		output, err := client.StartDBCluster(input)
		if err != nil {
			fmt.Println("Error while starting cluster ", *db)
			fmt.Println(err)
			fmt.Println(output)
		}
	}
}