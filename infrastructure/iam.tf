resource "aws_iam_role" "reddit_scraper_lambda_role" {
  name = "reddit-scraper-lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    }
  }]
}
EOF
}

resource "aws_iam_policy" "reddit_scraper_lambda_policy" {
  name        = "scraper-lambda-policy"
  description = "IAM policy for Lambda to access reddit SQS"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sqs:SendMessage",
        "sqs:SendMessageBatch",
        "sqs:GetQueueAttributes"
      ],
      "Resource": "${aws_sqs_queue.data_queue.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "reddit_scraper_lambda_basic_execution" {
  role       = aws_iam_role.reddit_scraper_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "reddit_scraper_lambda_policy_attachment" {
  role       = aws_iam_role.reddit_scraper_lambda_role.name
  policy_arn = aws_iam_policy.reddit_scraper_lambda_policy.arn
}

# Loader lambda
resource "aws_iam_role" "reddit_loader_lambda_role" {
  name = "reddit-loader-lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    }
  }]
}
EOF
}

resource "aws_iam_policy" "reddit_loader_lambda_policy" {
  name        = "loader-lambda-sqs-policy"
  description = "IAM policy for Loader Lambda to access SQS and dynamodb"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": [
        "${aws_dynamodb_table.posts_table.arn}",
        "${aws_dynamodb_table.comments_table.arn}"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
      ],
      "Resource": "${aws_sqs_queue.data_queue.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "reddit_loader_lambda_basic_execution" {
  role       = aws_iam_role.reddit_loader_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "reddit_loader_lambda_policy_attachment" {
  role       = aws_iam_role.reddit_loader_lambda_role.name
  policy_arn = aws_iam_policy.reddit_loader_lambda_policy.arn
}

# API lambda
resource "aws_iam_role" "reddit_api_lambda_role" {
  name = "reddit-api-lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    }
  }]
}
EOF
}

resource "aws_iam_policy" "reddit_api_lambda_policy" {
  name        = "api-lambda-sqs-policy"
  description = "IAM policy for API Lambda to access dynamodb"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:Query"
      ],
      "Resource": [
        "${aws_dynamodb_table.posts_table.arn}",
        "${aws_dynamodb_table.comments_table.arn}",
        "${aws_dynamodb_table.posts_table.arn}/index/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "reddit_api_lambda_basic_execution" {
  role       = aws_iam_role.reddit_api_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "reddit_api_lambda_policy_attachment" {
  role       = aws_iam_role.reddit_api_lambda_role.name
  policy_arn = aws_iam_policy.reddit_api_lambda_policy.arn
}

resource "aws_iam_role" "api_gateway_role" {
  name = "api-gateway-cloudwatch-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "apigateway.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "cloudwatch_logs_policy" {
  name = "api-gateway-cloudwatch-logs-policy"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      Effect   = "Allow",
      Resource = "*"
    }]
  })
}

resource "aws_iam_policy_attachment" "cloudwatch_logs_attachment" {
  name       = "api-gateway-cloudwatch-logs-attachment"
  roles      = [aws_iam_role.api_gateway_role.name]
  policy_arn = aws_iam_policy.cloudwatch_logs_policy.arn
}
