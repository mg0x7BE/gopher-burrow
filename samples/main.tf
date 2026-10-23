terraform {
  required_version = ">= 1.9"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "environment" {
  type        = string
  default     = "staging"
  description = "Deployment environment"
}

locals {
  name = "inventory-${var.environment}"
  tags = { Project = "inventory", Managed = true }
}

# Static assets for the web frontend
resource "aws_s3_bucket" "assets" {
  bucket = "${local.name}-assets"
  tags   = local.tags
}

resource "aws_instance" "api" {
  count         = var.environment == "prod" ? 3 : 1
  ami           = data.aws_ami.al2023.id
  instance_type = "t4g.small"

  lifecycle {
    create_before_destroy = true
  }
}

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]
}

output "bucket_arn" {
  value = aws_s3_bucket.assets.arn
}
