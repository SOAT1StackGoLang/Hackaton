resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "rds_subnet_group"
  subnet_ids = var.database_subnetids
}

resource "aws_security_group" "rds_security_group" {
  name        = "rds-security-group"
  vpc_id      = var.vpc_id
  description = "Allow all inbound for Postgress"

  ingress {
    from_port   = var.database_port
    to_port     = var.database_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_parameter_group" "postgres15_parameter_group" {
  name   = "postgres15-param-group"
  family = "postgres15"

    parameter {
      name  = "rds.force_ssl"
      value = "1"
    }
    parameter {
      name  = "max_connections"
      value = "1000"
    }
}

# Create KMS KEY to encrypt RDS

data "aws_caller_identity" "current" {}

resource "aws_kms_key" "rds_key" {
  description             = "KMS key for RDS instance encryption"
  deletion_window_in_days = 7
  enable_key_rotation     = true
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "key-default-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_kms_alias" "rds_key_alias" {
  name          = "alias/rds_key"
  target_key_id = aws_kms_key.rds_key.key_id
}

# Create RDS instance
resource "aws_db_instance" "rds" {
  identifier                 = "hackaton-rds-15"
  db_name                    = var.database_name
  instance_class             = "db.t4g.micro"
  allocated_storage          = 20
  storage_type               = "gp2"
  engine                     = "postgres"
  engine_version             = "15.6"
  auto_minor_version_upgrade = true
  skip_final_snapshot        = true
  publicly_accessible        = false
  vpc_security_group_ids     = [aws_security_group.rds_security_group.id]
  iam_database_authentication_enabled = false
  username                   = var.database_username
  password                   = var.database_password
  port                       = var.database_port
  db_subnet_group_name       = aws_db_subnet_group.rds_subnet_group.id
  availability_zone          = var.availability_zone
  multi_az                   = false
  parameter_group_name       = aws_db_parameter_group.postgres15_parameter_group.name

  # Enable encryption
  kms_key_id                 = aws_kms_key.rds_key.arn
  storage_encrypted = true

  ## add backup
  backup_retention_period = 7
}