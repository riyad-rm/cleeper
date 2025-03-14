package awsHandler

const FULL = 0
const START = 1
const STOP = 2
const IGNORE = 3

type LambdaTrigger struct {
        Action string `json:"action"`
        Regions string `json:"regions"`
        TagKeys string `json:"tagKeys"`
        TagValues string `json:"tagValues"`
        TaggedOnly string `json:"taggedOnly"` // using string instead of bool because if no value is provided by the user this will be set to false
}

type Config struct {
	// Services handling
	ec2 int
	rdsInstance int
	rdsCluster int
	asg int
	sagemaker int
	ecs int
	region []string
	tag_keys []string
	tag_values []string
}