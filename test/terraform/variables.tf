variable "region" {
  default = "eu-west-1"
}

variable "tags" {
  description = "common tags"
  default = {
    cleeper = "val1"
  }
}

variable "enable_rds_cluster"{
  default = 0
}