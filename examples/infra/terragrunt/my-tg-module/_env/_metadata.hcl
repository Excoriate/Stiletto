locals {
  tags = {
    ManagedBy      = "Terraform"
    OrchestratedBy = "Terragrunt"
    ModifiedAt     = run_cmd("sh", "-c", "export TERRAGRUNT_CFG_METADATA_RUNTIME_MODIFIEDAT=$(date +%Y-%m-%d); echo $TERRAGRUNT_CFG_METADATA_RUNTIME_MODIFIEDAT")
    Owner          = "MyCompanyOrOrg"
    Type           = "infrastructure"
  }
}
