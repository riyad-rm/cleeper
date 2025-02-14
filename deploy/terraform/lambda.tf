resource "aws_cloudwatch_log_group" "cleeper_log_group" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 14
}

resource "aws_lambda_function" "cleeper_lambda" {
  filename      = var.lambda_path
  function_name = var.lambda_function_name
  role          = aws_iam_role.cleeper_iam_for_lambda.arn
  source_code_hash = filebase64sha256(var.lambda_path)
  handler = var.lambda_function_name

  timeout = 600

  runtime = "go1.x"

  environment {
    variables = {
      version = "0.01"
    }
  }
}