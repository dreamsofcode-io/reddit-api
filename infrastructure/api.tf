resource "null_resource" "api_dependencies" {
  provisioner "local-exec" {
    command = "cd ${path.module}/../api && make build"
  }

  triggers = {
    index  = sha256(file("${path.module}/../api/go.sum"))
    go     = sha256(join("",fileset(path.module, "../api/*.go")))
    shared = sha256(join("",fileset(path.module, "../shared/**/*.go")))
  }
}

data "archive_file" "api_zip" {
  type        = "zip"
  source_file = "${path.module}/../bin/api"
  output_path = "${path.module}/.api.zip"

  depends_on = [
    resource.null_resource.api_dependencies
  ]
}

resource "aws_lambda_function" "api" {
  function_name    = "reddit-api"
  s3_bucket        = aws_s3_bucket.reddit-api-binary-bucket.id
  s3_key           = aws_s3_object.api_upload.key
  runtime          = "go1.x"
  handler          = "api"
  role             = aws_iam_role.reddit_api_lambda_role.arn
  memory_size      = 128
  timeout          = 30
  source_code_hash = data.archive_file.api_zip.output_base64sha256

  environment  {
    variables = {
      POST_TABLE_NAME = aws_dynamodb_table.posts_table.name
      COMMENT_TABLE_NAME = aws_dynamodb_table.comments_table.name
      POST_INDEX_NAME = "subreddit-timestamp-index"
    }
  }

  tags = local.tags
}

resource "aws_s3_object" "api_upload" {
  bucket = aws_s3_bucket.reddit-api-binary-bucket.id
  key    = "api.zip"
  source = data.archive_file.api_zip.output_path
  etag   = data.archive_file.api_zip.output_base64sha256
}

resource "aws_api_gateway_rest_api" "reddit" {
  name        = "reddit-api"
  description = "Reddit API Gateway"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = "${aws_api_gateway_rest_api.reddit.id}"
  parent_id   = "${aws_api_gateway_rest_api.reddit.root_resource_id}"
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = "${aws_api_gateway_rest_api.reddit.id}"
  resource_id   = "${aws_api_gateway_resource.proxy.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = "${aws_api_gateway_rest_api.reddit.id}"
  resource_id = "${aws_api_gateway_method.proxy.resource_id}"
  http_method = "${aws_api_gateway_method.proxy.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.api.invoke_arn}"
}

resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id   = "${aws_api_gateway_rest_api.reddit.id}"
  resource_id   = "${aws_api_gateway_rest_api.reddit.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = "${aws_api_gateway_rest_api.reddit.id}"
  resource_id = "${aws_api_gateway_method.proxy_root.resource_id}"
  http_method = "${aws_api_gateway_method.proxy_root.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.api.invoke_arn}"
}

resource "aws_api_gateway_deployment" "reddit" {
  depends_on = [
    aws_api_gateway_integration.lambda,
    aws_api_gateway_integration.lambda_root,
  ]

  rest_api_id = "${aws_api_gateway_rest_api.reddit.id}"
  stage_name  = "prod"
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:invokefunction"
  function_name = "${aws_lambda_function.api.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.reddit.execution_arn}/*/*"
}

resource "aws_lambda_permission" "apigw2" {
  statement_id  = "AllowAPIGatewayInvokeMore"
  action        = "lambda:invokefunction"
  function_name = "${aws_lambda_function.api.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.reddit.execution_arn}/*/*/*"
}
