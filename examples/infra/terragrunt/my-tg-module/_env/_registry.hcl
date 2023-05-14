locals {
  registry_base_url = get_env("TF_REGISTRY_BASE_URL", "git::https://github.com")
  registry_github_org= get_env("TF_REGISTRY_GITHUB_ORG", "Excoriate")
}
