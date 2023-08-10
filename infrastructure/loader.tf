resource "null_resource" "loader_lambda_dependencies" {
  provisioner "local-exec" {
    command = "cd ${path.module}/../loader && make build"
  }

  triggers = {
    index = sha256(file("${path.module}/../loader/go.sum"))
    go = sha256(join("",fileset(path.module, "../loader/*.go")))
  }
}

data "archive_file" "loader_zip" {
  type        = "zip"
  source_file = "${path.module}/../bin/loader"
  output_path = "${path.module}/.loader.zip"

  depends_on = [
    resource.null_resource.loader_lambda_dependencies
  ]
}

resource "aws_lambda_function" "loader_lambda" {
  function_name    = "reddit-loader"
  s3_bucket        = aws_s3_bucket.reddit-api-binary-bucket.id
  s3_key           = aws_s3_object.loader_upload.key
  runtime          = "go1.x"
  handler          = "loader"
  role             = aws_iam_role.reddit_loader_lambda_role.arn
  memory_size      = 128
  timeout          = 30
  source_code_hash = data.archive_file.loader_zip.output_base64sha256

  environment  {
    variables = {
      POST_TABLE_NAME = aws_dynamodb_table.posts_table.name
      COMMENT_TABLE_NAME = aws_dynamodb_table.comments_table.name
    }
  }

  tags = local.tags
}

resource "aws_lambda_event_source_mapping" "loader" {
  event_source_arn = aws_sqs_queue.data_queue.arn
  function_name = aws_lambda_function.loader_lambda.function_name
  batch_size = 10
}

resource "aws_s3_object" "loader_upload" {
  bucket = aws_s3_bucket.reddit-api-binary-bucket.id
  key    = "loader.zip"
  source = data.archive_file.loader_zip.output_path
  etag   = data.archive_file.loader_zip.output_base64sha256
}
