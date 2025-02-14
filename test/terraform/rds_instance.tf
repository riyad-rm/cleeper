resource "aws_db_instance" "default" {
  count = 0
  allocated_storage    = 10
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t3.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "foobarbaz65s4df6s4flmkkqklqa"
  parameter_group_name = "default.mysql5.7"
  skip_final_snapshot  = true

  tags = var.tags
}