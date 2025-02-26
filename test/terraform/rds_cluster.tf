resource "aws_rds_cluster" "default" {
  count = var.enable_rds_cluster
  cluster_identifier  = "aurora-cluster-demo"
  engine              = "aurora-mysql"
  engine_version      = "5.7.mysql_aurora.2.11.2"
  availability_zones  = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]
  database_name       = "mydb"
  master_username     = "foo"
  master_password     = "barqsdqd4q6d6q4f6q4daz4e1c61"
  skip_final_snapshot = true
  apply_immediately   = true
  tags                = var.tags
}


resource "aws_rds_cluster" "aurora1" {
  count = var.enable_rds_cluster
  cluster_identifier  = "aurora1-cluster-demo"
  availability_zones  = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]
  database_name       = "mydb"
  master_username     = "foo"
  master_password     = "barqsdqd4q6d6q4f6q4daz4e1c61"
  skip_final_snapshot = true
  apply_immediately   = true
  tags                = var.tags
}


resource "aws_rds_cluster" "postgresql" {
  count = var.enable_rds_cluster
  cluster_identifier  = "aurora-postgres-cluster-demo"
  engine              = "aurora-postgresql"
  availability_zones  = ["eu-west-1a", "eu-west-1b", "eu-west-1c"]
  database_name       = "mydb"
  master_username     = "foo"
  master_password     = "barqsdqd4q6d6q4f6q4daz4e1c61"
  skip_final_snapshot = true
  apply_immediately   = true
  tags                = var.tags
}


resource "aws_db_subnet_group" "default" {
  name       = "main"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "My DB subnet group"
  }
}

resource "aws_rds_cluster" "example" {
  count = var.enable_rds_cluster
  cluster_identifier   = "example"
  db_subnet_group_name = aws_db_subnet_group.default.name
  engine_mode          = "multimaster"
  master_password      = "barqsdqd4q6d6q4f6q4daz4e1c61"
  master_username      = "foo"
  skip_final_snapshot  = true
  apply_immediately    = true
  tags                 = var.tags
}


resource "aws_rds_global_cluster" "example" {
  count = var.enable_rds_cluster
  global_cluster_identifier = "global-test"
  engine                    = "aurora"
  engine_version            = "5.6.mysql_aurora.1.22.2"
  database_name             = "example_db"
}

resource "aws_rds_cluster" "primary" {
  count = var.enable_rds_cluster
  engine                    = aws_rds_global_cluster.example[0].engine
  engine_version            = aws_rds_global_cluster.example[0].engine_version
  cluster_identifier        = "test-primary-cluster"
  master_username           = "username"
  master_password           = "somepass123"
  database_name             = "example_db"
  global_cluster_identifier = aws_rds_global_cluster.example[0].id
  db_subnet_group_name      = "default"
  apply_immediately         = true
  skip_final_snapshot       = true
  tags                      = var.tags
}