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

data "aws_region" "current" {
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

  tags = {
    Table = "tibia-guild-members"
  }
}

resource "aws_dynamodb_table" "guilds_table" {
  name         = "tibia-guilds"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "date"

  attribute {
    name = "date"
    type = "S"
  }

  tags = {
    Table = "tibia-guilds"
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
        "Resource" : "arn:aws:logs:${data.aws_region.current.name}:*:*"
      }
    ]
  })
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
          "Resource" : "${aws_dynamodb_table.guild_member_action_table.arn}/*"
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
  function_name    = "get-tibia-guild-members-history"
  filename         = data.archive_file.get_guild_members_history.output_path
  source_code_hash = data.archive_file.get_guild_members_history.output_base64sha256

  role    = aws_iam_role.get_guild_members_history_role.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      GUILD_MEMBER_ACTION_TABLE = aws_dynamodb_table.guild_member_action_table.name
    }
  }
}

resource "aws_iam_role" "list_guilds_role" {
  name               = "list_guilds_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "list_guilds_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:Query",
            "dynamodb:Scan"
          ]
          "Resource" : aws_dynamodb_table.guilds_table.arn
        }
      ]
    })
  }
}

data "archive_file" "list_guilds" {
  type        = "zip"
  source_file = "functions/listguilds/main"
  output_path = "list_guilds.zip"
}

resource "aws_lambda_function" "list_guilds" {
  function_name    = "list-guilds"
  filename         = data.archive_file.list_guilds.output_path
  source_code_hash = data.archive_file.list_guilds.output_base64sha256

  role    = aws_iam_role.list_guilds_role.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      TIBIA_GUILDS_TABLE = aws_dynamodb_table.guilds_table.name
    }
  }
}

resource "aws_iam_role" "list_guilds_deaths_role" {
  name               = "list_guilds_deaths_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "list_guilds_inline_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:Query",
            "dynamodb:Scan",
          ]
          "Resource" : [
            aws_dynamodb_table.death_table.arn,
            "${aws_dynamodb_table.death_table.arn}/*",
          ]
        }
      ]
    })
  }
}

data "archive_file" "list_guilds_deaths" {
  type        = "zip"
  source_file = "functions/guilddeaths/main"
  output_path = "list_guild_deaths.zip"
}

resource "aws_lambda_function" "list_guilds_deaths" {
  function_name    = "list-guilds-deaths"
  filename         = data.archive_file.list_guilds_deaths.output_path
  source_code_hash = data.archive_file.list_guilds_deaths.output_base64sha256

  role    = aws_iam_role.list_guilds_deaths_role.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      DEATH_TABLE_NAME = aws_dynamodb_table.death_table.name
      DEATH_TABLE_CHARACTER_NAME_DATE_INDEX = "characterName-time-index"
      DEATH_TABLE_GUILD_TIME_INDEX= "guild-time-index"
    }
  }
}

resource "aws_apigatewayv2_api" "tibia" {
  name          = "tibia"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    allow_headers = ["Content-Type", "Authorization", "X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token"]
  }
}

resource "aws_apigatewayv2_stage" "tibia" {
  api_id      = aws_apigatewayv2_api.tibia.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "list_guilds_deaths" {
  api_id = aws_apigatewayv2_api.tibia.id

  integration_uri        = aws_lambda_function.list_guilds_deaths.invoke_arn
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "list_guilds_deaths" {
  api_id = aws_apigatewayv2_api.tibia.id

  route_key = "GET /guildDeaths/{guildName}"
  target    = "integrations/${aws_apigatewayv2_integration.list_guilds_deaths.id}"
}

resource "aws_lambda_permission" "api_gateway_list_guilds_deaths" {
  statement_id  = "api_gateway_list_guilds_deaths"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.list_guilds_deaths.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.tibia.execution_arn}/*/*"
}

resource "aws_apigatewayv2_integration" "get_guild_members_history" {
  api_id = aws_apigatewayv2_api.tibia.id

  integration_uri        = aws_lambda_function.get_tibia_guild_members_history.invoke_arn
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_guild_members_history" {
  api_id = aws_apigatewayv2_api.tibia.id

  route_key = "GET /guildMembersHistory/{guildName}"
  target    = "integrations/${aws_apigatewayv2_integration.get_guild_members_history.id}"
}

resource "aws_lambda_permission" "api_gateway_get_guild_members_history" {
  statement_id  = "api_gateway_get_guild_members_history"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_tibia_guild_members_history.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.tibia.execution_arn}/*/*"
}

resource "aws_apigatewayv2_integration" "list_guilds" {
  api_id = aws_apigatewayv2_api.tibia.id

  integration_uri        = aws_lambda_function.list_guilds.invoke_arn
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "list_guilds" {
  api_id = aws_apigatewayv2_api.tibia.id

  route_key = "GET /guildNames"
  target    = "integrations/${aws_apigatewayv2_integration.list_guilds.id}"
}

resource "aws_lambda_permission" "list_guilds" {
  statement_id  = "search_guild"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.list_guilds.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.tibia.execution_arn}/*/*"
}

resource "aws_amplify_app" "tibia_stats_ui" {
  name       = "tibia-stats-ui"
  repository = "https://github.com/jpodeszwik/tibia-stats-ui"

  build_spec = <<-EOT
    version: 1
    frontend:
      phases:
        preBuild:
          commands:
            - npm ci
        build:
          commands:
            - npm run build
      artifacts:
        baseDirectory: /dist
        files:
          - '**/*'
      cache:
        paths:
          - node_modules/**/*
  EOT

  custom_rule {
    source = "/<*>"
    status = "404-200"
    target = "/index.html"
  }

  custom_rule {
    source = "/api/<*>"
    status = "200"
    target = "${aws_apigatewayv2_api.tibia.api_endpoint}/<*>"
  }
}

resource "aws_route53domains_registered_domain" "tibia_stats_domain" {
  domain_name = "tibiastats.org"
}

resource "aws_amplify_domain_association" "example" {
  app_id      = aws_amplify_app.tibia_stats_ui.id
  domain_name = aws_route53domains_registered_domain.tibia_stats_domain.domain_name

  sub_domain {
    branch_name = "master"
    prefix      = ""
  }

  sub_domain {
    branch_name = "master"
    prefix      = "www"
  }
}

resource "aws_dynamodb_table" "death_table" {
  name         = "tibia-death"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "characterName"
  range_key    = "time"

  attribute {
    name = "characterName"
    type = "S"
  }

  attribute {
    name = "time"
    type = "S"
  }

  attribute {
    name = "guild"
    type = "S"
  }

  local_secondary_index {
    name            = "characterName-time-index"
    range_key       = "time"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "guild-time-index"
    hash_key        = "guild"
    range_key       = "time"
    projection_type = "ALL"
  }

  tags = {
    Table = "tibia-death"
  }
}

resource "aws_iam_user" "death_tracker" {
  name = "death-tracker"
}

resource "aws_iam_access_key" "death_tracker" {
  user = aws_iam_user.death_tracker.name
}

data "aws_iam_policy_document" "allow_death_tracker_death_table" {
  statement {
    effect  = "Allow"
    actions = [
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:PutItem",
      "dynamodb:BatchWriteItem",
    ]
    resources = [
      aws_dynamodb_table.death_table.arn,
      "${aws_dynamodb_table.death_table.arn}/*",
      aws_dynamodb_table.guild_exp_table.arn,
      "${aws_dynamodb_table.guild_exp_table.arn}/*",
      aws_dynamodb_table.guilds_table.arn,
      "${aws_dynamodb_table.guilds_table.arn}/*",
      aws_dynamodb_table.guilds_table.arn,
      "${aws_dynamodb_table.guilds_table.arn}/*",
      aws_dynamodb_table.highscore_table.arn,
      "${aws_dynamodb_table.highscore_table.arn}/*",
      aws_dynamodb_table.guild_members_table.arn,
      "${aws_dynamodb_table.guild_members_table.arn}/*",
      aws_dynamodb_table.guild_member_action_table.arn,
      "${aws_dynamodb_table.guild_member_action_table.arn}/*",
    ]
  }
}

resource "aws_iam_user_policy" "allow_death_tracker_death_table" {
  name   = "allow_death_tracker_death_table"
  user   = aws_iam_user.death_tracker.name
  policy = data.aws_iam_policy_document.allow_death_tracker_death_table.json
}

resource "aws_dynamodb_table" "guild_exp_table" {
  name         = "guild-exp"
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

  local_secondary_index {
    name            = "guildName-date-index"
    range_key       = "date"
    projection_type = "ALL"
  }

  tags = {
    Table = "tibia-guild-exp"
  }
}

resource "aws_iam_role" "get_guild_exp_role" {
  name               = "get_guild_exp_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  managed_policy_arns = [
    aws_iam_policy.lambda_log_policy.arn
  ]

  inline_policy {
    name   = "get_guild_exp_policy"
    policy = jsonencode({
      Version   = "2012-10-17",
      Statement = [
        {
          "Effect" : "Allow",
          "Action" : [
            "dynamodb:Query",
            "dynamodb:Scan",
          ]
          "Resource" : [
            aws_dynamodb_table.guild_exp_table.arn,
            "${aws_dynamodb_table.guild_exp_table.arn}/*",
          ]
        }
      ]
    })
  }
}

data "archive_file" "get_guild_exp" {
  type        = "zip"
  source_file = "functions/getguildexp/main"
  output_path = "get_guild_exp.zip"
}

resource "aws_lambda_function" "get_guild_exp" {
  function_name    = "get-guild-exp"
  filename         = data.archive_file.get_guild_exp.output_path
  source_code_hash = data.archive_file.get_guild_exp.output_base64sha256

  role    = aws_iam_role.get_guild_exp_role.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      GUILD_EXP_TABLE_NAME = aws_dynamodb_table.guild_exp_table.name
      GUILD_EXP_GUILD_NAME_DATE_INDEX = "guildName-date-index"
    }
  }
}

resource "aws_apigatewayv2_integration" "get_guild_exp" {
  api_id = aws_apigatewayv2_api.tibia.id

  integration_uri        = aws_lambda_function.get_guild_exp.invoke_arn
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_guild_exp" {
  api_id = aws_apigatewayv2_api.tibia.id

  route_key = "GET /guildExp/{guildName}"
  target    = "integrations/${aws_apigatewayv2_integration.get_guild_exp.id}"
}

resource "aws_lambda_permission" "get_guild_exp" {
  statement_id  = "get_guild_exp"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_guild_exp.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.tibia.execution_arn}/*/*"
}

resource "aws_dynamodb_table" "highscore_table" {
  name         = "tibia-highscore"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "worldName"
  range_key    = "date"

  attribute {
    name = "worldName"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  tags = {
    Table = "tibia-highscore"
  }
}

resource "aws_dynamodb_table" "guild_member_action_table" {
  name         = "tibia-guild-member-action"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "guildName"
  range_key    = "time-characterName"

  attribute {
    name = "guildName"
    type = "S"
  }

  attribute {
    name = "time-characterName"
    type = "S"
  }

  local_secondary_index {
    name            = "guildName-time-characterName-index"
    range_key       = "time-characterName"
    projection_type = "ALL"
  }

  tags = {
    Table = "tibia-guild-member-action"
  }
}
