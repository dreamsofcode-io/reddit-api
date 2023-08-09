resource "null_resource" "lambda_dependencies" {
  provisioner "local-exec" {
    command = "cd ${path.module}/../scraper && npm install"
  }

  triggers = {
    index = sha256(file("${path.module}/../scraper/index.js"))
    package = sha256(file("${path.module}/../scraper/package.json"))
    lock = sha256(file("${path.module}/../scraper/package-lock.json"))
    node = sha256(join("",fileset(path.module, "../scraper/**/*.js")))
  }
}

data "archive_file" "scraper_zip" {
  type        = "zip"
  source_dir  = "${path.module}/../scraper"
  output_path = "${path.module}/.lambda-bundle.zip"

  depends_on = [
    resource.null_resource.lambda_dependencies
  ]
}

resource "aws_lambda_function" "scraper_lambda" {
  function_name    = "reddit-scraper"
  filename         = data.archive_file.scraper_zip.output_path
  runtime          = "nodejs18.x"
  handler          = "index.handler"
  role             = aws_iam_role.reddit_scraper_lambda_role.arn
  memory_size      = 128
  timeout          = 600

  environment  {
    variables = {
      QUEUE_URL = aws_sqs_queue.data_queue.id
      CONNECTION_URL = var.connection_url
    }
  }

  tags = local.tags
}

resource "aws_cloudwatch_event_rule" "every_hour" {
    name = "every-hour"
    description = "Fires every hour"
    schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "run_scraper" {
    rule = aws_cloudwatch_event_rule.every_hour.name
    target_id = "run_scraper"
    arn = aws_lambda_function.scraper_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_scraper" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.scraper_lambda.function_name
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.every_hour.arn
}