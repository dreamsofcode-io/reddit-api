resource "random_pet" "lambda_binary_bucket_name" {
  prefix = "dreamsofcode"
  length = 2
}

resource "aws_s3_bucket" "reddit-api-binary-bucket" {
  bucket = "reddit-api-${random_pet.lambda_binary_bucket_name.id}"

  tags = local.tags
}
