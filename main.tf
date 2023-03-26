terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "eu-central-1"

  default_tags {
    tags = {
      Project = "tibia-data"
    }
  }
}

resource "aws_dynamodb_table" "exp_table" {
  name         = "tibia-exp"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "playerName"
  range_key    = "date"

  attribute {
    name = "playerName"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  global_secondary_index {
    name            = "playerName-date-index"
    hash_key        = "playerName"
    range_key       = "date"
    projection_type = "ALL"
  }
}

resource "aws_dynamodb_table" "guild_members_table" {
  name         = "tibia-guild-members"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "guildName"
  range_key    = "date"

  attribute {
    name = "guildName"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  global_secondary_index {
    name            = "guildName-date-index"
    hash_key        = "guildName"
    range_key       = "date"
    projection_type = "ALL"
  }
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_policy" "lambda_log_policy" {
  name   = "log_stream_policy"
  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        "Effect" : "Allow",
        "Action" : [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:CreateLogGroup"
        ]
        "Resource" : "arn:aws:logs:eu-central-1:*:*"
      }
    ]
  })
}

resource "aws_iam_role" "get_tibia_exp_role" {
  name               = "get_tibia_exp_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "get_tibia_exp_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:Query",
            "dynamodb:Scan"
          ]
          "Resource" : "${aws_dynamodb_table.exp_table.arn}/*"
        }
      ]
    })
  }
}

data "archive_file" "get_exp" {
  type        = "zip"
  source_file = "functions/getexp/main"
  output_path = "get_exp.zip"
}

resource "aws_lambda_function" "get_tibia_exp" {
  function_name = "get-tibia-exp"
  filename      = "get_exp.zip"
  role          = aws_iam_role.get_tibia_exp_role.arn
  handler       = "main"

  source_code_hash = data.archive_file.get_exp.output_base64sha256

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      TIBIA_EXP_TABLE = aws_dynamodb_table.exp_table.name
    }
  }
}

resource "aws_iam_role" "get_guild_members_history_role" {
  name               = "get_guild_members_history_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "get_guild_members_history_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:Query",
            "dynamodb:Scan"
          ]
          "Resource" : "${aws_dynamodb_table.guild_members_table.arn}/*"
        }
      ]
    })
  }
}

data "archive_file" "get_guild_members_history" {
  type        = "zip"
  source_file = "functions/guildhistory/main"
  output_path = "get_guild_members_history.zip"
}

resource "aws_lambda_function" "get_tibia_guild_members_history" {
  function_name = "get-tibia-guild-members-history"
  filename      = "get_guild_members_history.zip"
  role          = aws_iam_role.get_guild_members_history_role.arn
  handler       = "main"

  source_code_hash = data.archive_file.get_exp.output_base64sha256

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      TIBIA_GUILD_MEMBERS_TABLE = aws_dynamodb_table.guild_members_table.name
    }
  }
}

resource "aws_iam_role" "load_players_exp" {
  name               = "load_players_exp"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "load_players_exp_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:BatchWriteItem",
            "dynamodb:PutItem"
          ]
          "Resource" : aws_dynamodb_table.exp_table.arn
        }
      ]
    })
  }
}

data "archive_file" "load_players_exp" {
  type        = "zip"
  source_file = "functions/exp/main"
  output_path = "load_players_exp.zip"
}

resource "aws_lambda_function" "load_players_exp" {
  function_name = "load-players-exp"
  filename      = "load_players_exp.zip"
  role          = aws_iam_role.load_players_exp.arn
  handler       = "main"

  source_code_hash = data.archive_file.get_exp.output_base64sha256

  runtime = "go1.x"

  timeout = 300

  environment {
    variables = {
      TIBIA_EXP_TABLE = aws_dynamodb_table.exp_table.name
    }
  }
}

resource "aws_iam_role" "load_guild_members" {
  name               = "load_guild_members"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "load_guild_members_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:BatchWriteItem",
            "dynamodb:PutItem"
          ]
          "Resource" : aws_dynamodb_table.guild_members_table.arn
        }
      ]
    })
  }
}

data "archive_file" "load_guild_members" {
  type        = "zip"
  source_file = "functions/guild/main"
  output_path = "load_guild_members.zip"
}

resource "aws_lambda_function" "load_guild_members" {
  function_name = "load-tibia-guild-members"
  filename      = "load_guild_members.zip"
  role          = aws_iam_role.load_guild_members.arn
  handler       = "main"

  source_code_hash = data.archive_file.get_exp.output_base64sha256

  runtime = "go1.x"

  timeout = 300

  environment {
    variables = {
      TIBIA_GUILD_MEMBERS_TABLE = aws_dynamodb_table.guild_members_table.name
    }
  }
}
