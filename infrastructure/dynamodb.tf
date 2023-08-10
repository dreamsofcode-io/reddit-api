resource "aws_dynamodb_table" "posts_table" {
  name           = "reddit-posts-table"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "post_id"

  global_secondary_index {
    hash_key           = "subreddit"
    name               = "subreddit-timestamp-index"
    projection_type    = "ALL"
    range_key          = "timestamp"
  }

  attribute {
    name = "post_id"
    type = "S"
  }

  attribute {
    name = "subreddit"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  tags = local.tags
}

resource "aws_dynamodb_table" "comments_table" {
  name           = "reddit-comments-table"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "post_id"

  attribute {
    name = "post_id"
    type = "S"
  }

  tags = local.tags
}
