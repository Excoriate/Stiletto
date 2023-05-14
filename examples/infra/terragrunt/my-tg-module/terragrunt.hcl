locals {
  aws_region                               = get_env("TF_VAR_aws_region", "us-east-1")
  environment                          = get_env("TF_VAR_environment")
  /*
    NOTE:
    ----------------------------------------------------------
    * Customise state management for this root module
    ----------------------------------------------------------
  */
  terraform_state_file_bucket_region       = get_env("TF_STATE_BUCKET_REGION", local.aws_region)
  terraform_state_file_bucket              = get_env("TF_STATE_BUCKET")
  terraform_state_file_lock_dynamodb_table = get_env("TF_STATE_LOCK_TABLE")

  /*
    NOTE:
    ----------------------------------------------------------
    * Customise binary versions accordingly.
    ----------------------------------------------------------
  */
  terraform_version = get_env("TF_VERSION", "v1.4.6")
  terragrunt_version = get_env("TG_VERSION", "v0.42.8")
}

terraform {
  extra_arguments "optional_vars" {
    commands = [
      "apply",
      "destroy",
      "plan",
    ]

    required_var_files = [
      "${get_terragrunt_dir()}/../config/common.tfvars",
      "${get_terragrunt_dir()}/../config/common-dev.tfvars",
      "${get_terragrunt_dir()}/../config/common-prod.tfvars",
    ]

    optional_var_files = [
      "${local.aws_region}.tfvars", // Global values that apply to all environments for a certain region.
      "${local.environment}.tfvars", // Environment specific values that apply to all regions.
      "${local.environment}-${local.aws_region}.tfvars", // Environment and region specific values.
    ]
  }

  extra_arguments "disable_input" {
    commands  = get_terraform_commands_that_need_input()
    arguments = ["-input=false"]
  }

  after_hook "clean_cache_after_apply" {
    commands = ["apply"]
    execute  = ["rm", "-rf", ".terragrunt-cache"]
  }

  after_hook "remove_auto_generated_backend" {
    commands = ["apply"]
    execute  = ["rm", "-rf", "backend.tf"]
  }

  after_hook "remove_auto_generated_provider" {
    commands = ["apply"]
    execute  = ["rm", "-rf", "provider.tf"]
  }
}


generate "terraform_version" {
  path              = ".terraform-version"
  if_exists         = "overwrite"
  disable_signature = true

  contents = <<-EOF
    ${local.terraform_version}
  EOF
}

generate "terragrunt_version" {
  path              = ".terragrunt-version"
  if_exists         = "overwrite"
  disable_signature = true

  contents = <<-EOF
    ${local.terragrunt_version}
  EOF
}


generate "providers" {
  path      = "providers.tf"
  if_exists = "overwrite_terragrunt"

  contents = templatefile("${get_repo_root()}/examples/infra/terragrunt/my-tg-module/providers.tf.tmpl", {
    aws_region_passed_from_env = local.aws_region
  })
}

remote_state {
  backend = "s3"

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }

  config = {
    disable_bucket_update = true
    encrypt               = true

    region         = local.terraform_state_file_bucket_region
    dynamodb_table = local.terraform_state_file_lock_dynamodb_table
    bucket         = local.terraform_state_file_bucket

    key =  "rotation-lambda/${replace(path_relative_to_include(), "/^\\d+[_-]*/", "")}/${local.environment}/${local.aws_region}/terraform.tfstate"
  }
}
