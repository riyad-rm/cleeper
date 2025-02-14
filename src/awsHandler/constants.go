package awsHandler


// Specify max count possible per call
// In order to try and retrieve the maximum amount
// objects and avoid multiple NextToken calls

const MAX_EC2_INSTANCES = 50
const MAX_ASG_INSTANCES = 50
const MAX_RDS_INSTANCES = 100