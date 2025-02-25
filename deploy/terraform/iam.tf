resource "aws_iam_role" "cleeper_iam_for_lambda" {
  name = "cleeper_iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "cleeper_lambda_logging" {
  name        = "cleeper_lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:log-group:/aws/lambda/${var.lambda_function_name}:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cleeper_lambda_logs" {
  role       = aws_iam_role.cleeper_iam_for_lambda.name
  policy_arn = aws_iam_policy.cleeper_lambda_logging.arn
}

resource "aws_iam_policy" "cleeper_lambda_permissions" {
  name        = "cleeper_lambda_permissions"
  path        = "/"
  description = "IAM policy granting permissions to cleeper lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid" : "ASG",
      "Action": [
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:SuspendProcesses",
        "autoscaling:ResumeProcesses"
      ],
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Sid":"EC2",
      "Action": [
        "ec2:StartInstances",
        "ec2:TerminateInstances",
        "ec2:StopInstances",
        "ec2:DescribeRegions",
        "ec2:DescribeInstances"
      ],
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Sid":"RDS",
      "Action": [
        "rds:StartDBInstance",
        "rds:StopDBInstance",
        "rds:DescribeDBClusters",
        "rds:DescribeDBInstances",
        "rds:StopDBCluster",
        "rds:StartDBCluster"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cleeper_lambda_permissions" {
  role       = aws_iam_role.cleeper_iam_for_lambda.name
  policy_arn = aws_iam_policy.cleeper_lambda_permissions.arn
}