resource "aws_launch_template" "cleeper_template" {
  name_prefix   = "cleeper"
  image_id      = "ami-0088366b4b407a312"
  instance_type = "t2.micro"
}

resource "aws_autoscaling_group" "cleeper_asg" {
  availability_zones = ["eu-west-1a"]
  desired_capacity   = 0
  max_size           = 1
  min_size           = 0

  launch_template {
    id      = aws_launch_template.cleeper_template.id
    version = "$Latest"
  }

  tag {
    key                 = "cleeper"
    value               = "val1"
    propagate_at_launch = true
  }
}