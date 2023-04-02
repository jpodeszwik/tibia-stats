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

resource "aws_dynamodb_table" "guilds_table" {
  name         = "tibia-guilds"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "date"

  attribute {
    name = "date"
    type = "S"
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
  function_name    = "get-tibia-exp"
  filename         = data.archive_file.get_exp.output_path
  source_code_hash = data.archive_file.get_exp.output_base64sha256

  role    = aws_iam_role.get_tibia_exp_role.arn
  handler = "main"


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
  function_name    = "get-tibia-guild-members-history"
  filename         = data.archive_file.get_guild_members_history.output_path
  source_code_hash = data.archive_file.get_guild_members_history.output_base64sha256

  role    = aws_iam_role.get_guild_members_history_role.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 10

  environment {
    variables = {
      TIBIA_GUILD_MEMBERS_TABLE = aws_dynamodb_table.guild_members_table.name
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
  source_file = "functions/etlexp/main"
  output_path = "load_players_exp.zip"
}

resource "aws_lambda_function" "load_players_exp" {
  function_name    = "load-players-exp"
  filename         = data.archive_file.load_players_exp.output_path
  source_code_hash = data.archive_file.load_players_exp.output_base64sha256

  role    = aws_iam_role.load_players_exp.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 600

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
          "Resource" : [
            aws_dynamodb_table.guild_members_table.arn,
            aws_dynamodb_table.guilds_table.arn
          ]
        }
      ]
    })
  }
}

data "archive_file" "load_guild_members" {
  type        = "zip"
  source_file = "functions/etlguild/main"
  output_path = "load_guild_members.zip"
}

resource "aws_lambda_function" "load_guild_members" {
  function_name    = "load-tibia-guild-members"
  filename         = data.archive_file.load_guild_members.output_path
  source_code_hash = data.archive_file.load_guild_members.output_base64sha256

  role    = aws_iam_role.load_guild_members.arn
  handler = "main"

  runtime = "go1.x"

  timeout = 300

  environment {
    variables = {
      TIBIA_GUILD_MEMBERS_TABLE = aws_dynamodb_table.guild_members_table.name
      TIBIA_GUILDS_TABLE = aws_dynamodb_table.guilds_table.name
    }
  }
}

resource "aws_cloudwatch_event_rule" "load_player_exp" {
  name                = "load-player-exp-schedule"
  schedule_expression = "cron(0 11 * * ? *)"
}

resource "aws_cloudwatch_event_target" "load_player_exp" {
  rule      = aws_cloudwatch_event_rule.load_player_exp.name
  target_id = "lambda"
  arn       = aws_lambda_function.load_players_exp.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_load_player_exp" {
  statement_id  = "allow_cloudwatch_to_load_player_exp"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.load_players_exp.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.load_player_exp.arn
}

resource "aws_cloudwatch_event_rule" "load_guild_members" {
  name                = "load-guild-members-schedule"
  schedule_expression = "cron(0 11 * * ? *)"
}

resource "aws_cloudwatch_event_target" "load_guild_members" {
  rule      = aws_cloudwatch_event_rule.load_guild_members.name
  target_id = "lambda"
  arn       = aws_lambda_function.load_guild_members.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_load_guild_members" {
  statement_id  = "allow_cloudwatch_to_load_guild_members"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.load_guild_members.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.load_guild_members.arn
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
  api_id = aws_apigatewayv2_api.tibia.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "get_player_exp" {
  api_id = aws_apigatewayv2_api.tibia.id

  integration_uri        = aws_lambda_function.get_tibia_exp.invoke_arn
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "get_player_exp" {
  api_id = aws_apigatewayv2_api.tibia.id

  route_key = "GET /playerExp/{playerName}"
  target    = "integrations/${aws_apigatewayv2_integration.get_player_exp.id}"
}

resource "aws_lambda_permission" "api_gateway_get_player_exp" {
  statement_id  = "api_gateway_get_player_exp"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_tibia_exp.function_name
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