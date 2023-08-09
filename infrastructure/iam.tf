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
  name        = "lambda-sqs-policy"
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
