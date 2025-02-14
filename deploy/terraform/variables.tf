variable "lambda_function_name" {
  description = "Cleeper Lambda name."
  default = "cleeper"
}

variable "lambda_path" {
  type = string
  description = "path to lambda package"
  default = "../../build/bin/cleeper.zip"
}

variable "region" {
  type = string
  description = "region where to deploy"
  default = "eu-west-1"
}