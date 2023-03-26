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
  name           = "tibia-exp"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "playerName"
  range_key      = "date"

  attribute {
    name = "playerName"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  global_secondary_index {
    name               = "playerName-date-index"
    hash_key           = "playerName"
    range_key          = "date"
    projection_type    = "ALL"
  }
}

resource "aws_dynamodb_table" "guild_members_table" {
  name           = "tibia-guild-members"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "guildName"
  range_key      = "date"

  attribute {
    name = "guildName"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  global_secondary_index {
    name               = "guildName-date-index"
    hash_key           = "guildName"
    range_key          = "date"
    projection_type    = "ALL"
  }
}
