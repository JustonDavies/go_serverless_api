//-- Terraform Configuration -------------------------------------------------------------------------------------------

//-- Providers ---------------------------------------------------------------------------------------------------------
variable aws_region {}
variable aws_profile {}
provider "aws" {
  version    = "2.7"

  region  = "${var.aws_region}"
  profile = "${var.aws_profile}"
}

//-- Shared Internal Infrastructure ------------------------------------------------------------------------------------
module "shared_internal_vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> v2.0"

  name = "shared-internal-vpc"
  cidr = "10.0.0.0/16"

  azs = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]

  public_subnets      = ["10.0.10.0/24", "10.0.11.0/24", "10.0.12.0/24"]
  private_subnets     = ["10.0.20.0/24", "10.0.21.0/24", "10.0.22.0/24"]
  database_subnets    = ["10.0.30.0/24", "10.0.31.0/24", "10.0.32.0/24"]

  enable_nat_gateway = true

  tags = {
    terraform   = "true"
    application = "shared_internal"
  }
}

resource "aws_security_group" "lambda_security_group" {
  vpc_id      = "${module.shared_internal_vpc.vpc_id}"

  tags = {
    name        = "lambda_security_group"
    Name        = "lambda_security_group"
    terraform   = "true"
    application = "shared_internal"
  }
}

  resource "aws_security_group_rule" "lambda_security_group_allow_internet_egress_rule" {
    type = "egress"
    security_group_id = "${aws_security_group.lambda_security_group.id}"

    description = "Allow Lambdas full tcp access to the internet"

    protocol = "all"
    from_port = 0
    to_port = 65535

    cidr_blocks = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

resource "aws_security_group" "postgres_security_group" {
  vpc_id      = "${module.shared_internal_vpc.vpc_id}"

  tags = {
    name        = "postgres_security_group"
    Name        = "postgres_security_group"
    terraform   = "true"
    application = "shared_internal"
  }
}

  resource "aws_security_group_rule" "postgres_security_group_allow_lambda_security_group_rule" {
    type = "ingress"
    security_group_id = "${aws_security_group.postgres_security_group.id}"

    description = "Allow Lambdas protected access to the PostgreSQL"

    protocol = "tcp"
    from_port = 5432
    to_port = 5432

    source_security_group_id = "${aws_security_group.lambda_security_group.id}"
  }

//-- Application: Go serverless API ------------------------------------------------------------------------------------
variable go_serverless_api_database_password {}
resource "aws_db_instance" "go_serverless_api_database" {
  identifier = "go-serverless-api-database"

  engine            = "postgres"
  engine_version    = "11.1"

  instance_class    = "db.t2.micro"
  allocated_storage = 20

  name     = "go_serverless_api_database"
  username = "go_serverless_api_worker"
  password = "${var.go_serverless_api_database_password}"
  port     = "5432"

  vpc_security_group_ids = ["${aws_security_group.postgres_security_group.id}"]
  db_subnet_group_name   = "${module.shared_internal_vpc.database_subnet_group}"

  skip_final_snapshot = true
  final_snapshot_identifier = "go-serverless-api-database-snapshot"

  tags = {
    terraform   = "true"
    application = "task_service"
  }
}

variable stage{}
output "task_service_serverless_secrets" {
  value = <<EOF

#-- BEGIN TERRAFORM GENERATED SECRETS ----------
aws:
  stage: ${var.stage}
  region: ${var.aws_region}
  profile: ${var.aws_profile}

  rds:
    engine: ${aws_db_instance.go_serverless_api_database.engine}
    url: ${aws_db_instance.go_serverless_api_database.address}
    name: ${aws_db_instance.go_serverless_api_database.name}
    username: ${aws_db_instance.go_serverless_api_database.username}
    password: ${aws_db_instance.go_serverless_api_database.password}
    ssl_mode: require

  vpc:
   subnet_ids: [${join(",", module.shared_internal_vpc.private_subnets)}]
   security_group_ids: [${aws_security_group.lambda_security_group.id}]

  schedule:
    warming: cron(0/5 * * * ? *)

#-- END TERRAFORM GENERATED SECRETS ----------

  EOF
}

//-- Application: Something Else ---------------------------------------------------------------------------------------